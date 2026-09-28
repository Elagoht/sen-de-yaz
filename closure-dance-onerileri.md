# Collage: form action'larındaki "closure dansı" için çözüm önerileri

> Kapsam notu: Bu öneriler framework'ü tüketici tarafından okuyarak hazırlandı
> (sen-de-yaz uygulamasının collage v0.32.0'ı kullanım biçimi). Önerilen API
> adları (`rc.Page`, `WithReRender` vb.) temsilidir; `RenderContext` ve
> `ActionResult`'ın iç yapısını görmeden yazıldığı için ayrıntılarda
> kaymaya açık. Kaynak kodla karşılaştırınca gözden geçirmeye değer.

## 1. Teşhis: dans neden var

Her form sayfası bugün aynı dört adımı tekrarlıyor
(`pages/panel/stories/form.go`):

```go
func CreatePage(storyService *storydomain.StoryService, userService *users.UserService) *collage.Page {
	var page *collage.Page // 1. henüz yok
	action := collage.NewAction("story-create").
		WithMethods(http.MethodPost).
		WithHandler(actions.CreateStoryAction(storyService, userService,
			func() *collage.Page { return page }, // 2. sonra var olacak, söz veriyoruz
		)).Build()
	page = collage.NewPage("story-create"). // 3. sayfayı kur
		WithLayouts(layouts.Layout(), layouts.PanelLayout(userService)).
		WithContent(stories.StoryCreateBlock().WithDataHandler(collage.Load(createData)).Build()).
		WithPath("en", "/stories/new").
		WithActionFor(action). // 4. eylemi sayfaya bağla
		Build()
	return page
}
```

Ve her action, bu dolaylılığı imzasında taşıyor (`actions/stories.go`):

```go
func CreateStoryAction(
	service *stories.StoryService,
	userService *users.UserService,
	page func() *collage.Page, // tek bir dönüş yolu için tüm imzayı kirletiyor
) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		// ...
		if len(fieldErrors) > 0 {
			return utilities.FormErrors(page, rc, fieldErrors) // 422 + sayfayı yeniden render
		}
		return collage.SeeOther(fmt.Sprintf("/stories/%d", story.ID)), nil
	}
}
```

`utilities.FormErrors` da sayfa referansını sonuca işliyor:

```go
return &collage.ActionResult{Status: http.StatusUnprocessableEntity, Page: page()}, nil
```

Maliyet dört kat:

1. **Ritüel tekrarı.** Her form sayfası aynı `var page` + closure + `WithActionFor`
   koreografisini yazıyor.
2. **Görünmez kırılganlık.** Dizilim sırasına güveniyor: biri "temizlik" diye
   sayfayı üste alırsa handler artık nil sayfa yakalar; hata build'de değil,
   o form ilk kez 422 döndüğünde, production'da nil-deref olarak çıkar.
3. **Kirli imzalar.` CreateStoryAction`ın `page func() *collage.Page`
   parametresi, action'ın gerçekte bağımlı olmadığı bir şeye bağımlıymış gibi
   davranır — okuyucu yanlış yönlendirilir, test kurulumu şişer.
4. **İki ayrı dünya.** Sayfaya bağlı action'lar (`WithActionFor`) ve bağımsız
   kayıtlı action'lar (`routes.go`'daki `editEntry` gibi) aynı ihtiyacı farklı
   hack'lerle çözüyor; `editEntry` şanslı, çünkü `detail` ondan önce kuruluyor.

## 2. Ana gözlem

Dansın sebebi: handler **kurulum anında** yazılıyor ama ihtiyaç **istek
anında** ortaya çıkıyor — ve kurulum anında sayfa henüz yok.

Oysa istek anında framework sayfayı zaten biliyor:

- `WithActionFor(action)` ile bağlanan action, isteği yalnızca o sayfanın
  URL'sinden alabilir (action'ın kendi path'leri sayfayla değiştiriliyor).
- `WithAction("POST", h)` kısayolu action'ı zaten sayfanın içinde tanımlıyor.

Yani `ActionResult.Page`'i handler'ın söylemek zorunda olmasının tek istisnası,
`routes.go`'daki gibi path'i sayfadan farklı olan bağımsız action'lar. Bunun
dışında **framework'ün bilmediği bir şeyi handler'a söyletiyoruz.**

## 3. Öneri 1 — RenderContext sayfayı açsın (çekirdek değişiklik)

Dispatch, sayfaya bağlı bir action'ı çalıştırırken `RenderContext`'e sahibi
olan sayfayı koyar:

```go
// collage tarafından doldurulur
rc.Page *Page
```

Ve dispatch kuralı genişler: **action, re-render anlamı taşıyan bir status
(422) döndürürse ve `Page` nil ise, framework action'ın bağlı olduğu sayfayı
koyar.**

```go
// actions/stories.go — imza temizlenir
func CreateStoryAction(service *stories.StoryService, userService *users.UserService) collage.ActionHandlerFunc {
	return func(ctx context.Context, rc *collage.RenderContext) (*collage.ActionResult, error) {
		// ...
		if len(fieldErrors) > 0 {
			return &collage.ActionResult{Status: http.StatusUnprocessableEntity}, nil
			// Page nil → dispatch bağlı sayfayı doldurur
		}
		return collage.SeeOther(fmt.Sprintf("/stories/%d", story.ID)), nil
	}
}
```

Geriye dönük uyumlu: bugün `Page: page()` gönderen hiçbir kod bozulmaz
(açık değer fallback'i yener). Bugün nil `Page` ile 422 döndürmek muhtemelen
zaten bozuk davranış olduğu için, değişen tek şey "bozuk"un "çalışır" olması.

## 4. Öneri 2 — Doğrulama hatası birinci sınıf sonuç olsun

422 + alan hataları + yeniden render üçlüsü her framework'ün kullanıcısında
aynı olacak; collage bu üçlüyü sahiplensin:

```go
// RenderContext üzerinde şeker
func (rc *RenderContext) FormErrors(errors map[string]string) *ActionResult {
	rc.Set("form_errors", errors) // bugünkü rc.Set("form_errors", ...) sözleşmesi korunur
	return &ActionResult{Status: http.StatusUnprocessableEntity}
}
```

```go
if len(fieldErrors) > 0 {
	return rc.FormErrors(fieldErrors), nil
}
```

Bununla `utilities.FormErrors` tamamen kaybolur ve `FormView`/`FormData`
şablon tarafında kalır. İki yan kazanç:

- **form_errors artık framework'ün resmî sözleşmesi** olduğu için ileride
  `{{fieldError "title"}}` benzeri şablon yardımcıları tek yerde tanımlanabilir.
- Re-render kararı tek noktaya iner; ileride fragment-seviyesinde kısmi
  re-render (yalnızca formun olduğu `{{slot}}`) aynı musluktan verilebilir —
  `WithFragmentPath`/HTMX yönünde giden yolun zeminini hazırlar.

## 5. Öneri 3 — Bağımsız action'lar için `WithReRender(name)`

`routes.go`'daki `editEntry` kasıtlı olarak bağımsız: sayfası `/stories/{id}`
ama action `/stories/{id}/edit` üzerinden POST alıyor, `WithActionFor` ise
path'leri sayfayla değiştirdiği için bu action sayfaya bağlanamaz. Bu dünyanın
açık bir çıkışa ihtiyacı var:

```go
editEntry := collage.NewAction("story-edit-entry").
	WithPath("en", "/stories/{id}/edit").
	WithMethods(http.MethodPost).
	WithReRender("story-detail"). // isimle çöz; app başlarken doğrulanır
	WithHandler(actions.UpdateEntryAction(storyService, userService)).
	Build()
```

Doğrulama, framework'ün mevcut felsefesine oturur: builder hataları birikir,
**app start'ta** (tüm `register` çağrıları bittikten sonra) "story-detail adında
bir sayfa yok" hatası patlar. Sayfa adlarıyla startup'ta çözümleme collage'ta
zaten kabul görmüş bir deyim — şablonlar `layouts/default.html` diye path'le
referanslanıyor.

Statik olarak denetlenemeyen tek şey kalır: `WithReRender` beyan etmemiş,
sayfaya bağlı olmayan bir action çalışma anında nil-`Page`'li 422 döndürürse.
Bunun için ceremony değil, net mesaj: `action %q 422 döndürdü ama re-render
hedefi yok; sayfaya bağlayın ya da WithReRender(name) beyan edin`.

## 6. Değerlendirilen ve elenen alternatifler

| Alternatif | Neden elendi |
|---|---|
| `*PageRef` tanıtıcıları / DI token'ı (action ref yakalar, `RegisterPage` bağlar) | Yeni bir kavram ekler, Öneri 1'in verdiği sonucu verir. |
| Handler'a build anında `*Page` vermek | İmkânsız — sorun zaten sayfanın o anda var olmaması. (closure-thunk tam olarak bunun el yapımı hâli.) |
| Sadece mevcut deseni dokümante etmek | Sıralama tuzağı ve ritüel yerinde kalır. |
| `Page`'i build sonrası değiştirilebilir kılmak (two-phase build) | Framework'ün genel tarzı olan "kur → registration'da doğrula → dondur" çizgisini kırar. |

## 7. Göç yolu (üç küçük sürüm, her adım geriye dönük uyumlu)

1. **rc.Page + nil-Page fallback'i.** Hiçbir uygulama kodu bozulmaz; yeni
   action'lar closure'suz yazılabilir.
2. **rc.FormErrors / ValidationFailure.** `utilities.FormErrors` ince bir
   sarmalayıcı kalabilir ya da silinir; her sayfa tek tek sadeleşir.
   Bu adımdan sonra sen-de-yaz'daki hiçbir action `func() *collage.Page`
   almaz — `form.go`'daki dans biter.
3. **WithReRender + startup doğrulaması.** `routes.go`'daki `editEntry`
   closure'ı da gider; kalan tek "sır", sayfa path'i dışında path'i olan
   action'ların isimle hedef beyan etmesi — bu da artık belgelenmiş, denetlenen
   bir kural.

## 8. Kenar durumları

- **Aynı handler iki sayfaya bağlıysa:** `rc.Page` istek başına çözülür;
  ek iskele gerekmez. İsim temelli tek çözümde (Öneri 1'siz saf `WithReRender`)
  bu durum bozulur — bu yüzden çekirdek, bağlı-sayfa fallback'idir; isim yalnız
  bağımsız action'lar içindir.
- **Açık `Page` değeri** her zaman fallback'i yener (kaçış kapağı korunur).
- **Yönlendirmeler** (`SeeOther`) etkilenmez; fallback yalnız re-render
  anlamındaki status'lara (422) uygulanır, 200/3xx'e değil.
- **Honeypot/CSRF:** 422 re-render'ı bugünkü gibi tam sayfa render olduğundan
  plugin'ler token'ı yeniler; framework tarafına taşınınca da davranış aynı
  kalmalı (uygularken doğrulanmalı).
- **Testler:** handler artık sayfa kurmadan test edilebilir; rc hazırken
  `rc.Page`'e test sayfası takmak yeterli.
