## How to Run Application?
After cloning the project, follow the steps:
1. Run `docker compose up -d` 
#####
2. Run `migrate -source file://pkg/common/database/migrations -database "postgres://postgres:postgres@localhost:5435/cart_db?sslmode=disable" up 2` (**WARNING: run this command just once!**)


## How to Run Integration Tests?
1. Run `docker compose up -d` command if you not did not run already
#####
2. Run `migrate -source file://pkg/common/database/migrations -database "postgres://postgres:postgres@localhost:5436/cart_db_test?sslmode=disable" up 2` (**WARNING: run this command just once!**)
#####
3. Set environment variables as:
- `export ENVIRONMENT=TEST`
- `export TEST_DB_URL=postgres://postgres:postgres@localhost:5436/cart_db_test?sslmode=disable`
#####
4. Run `go test -run Test` in `./pkg/handlers/item/integration_tests` for integration tests in item package, and run in `./pkg/handlers/cart/integration_tests` for cart packages tests.

## How to Run Unit Tests?
1. Run `docker compose up -d` command if you not did not run already
#####
2. Run `go test -run Test` in  `./pkg/handlers/item` for items unit tests, and run in `./pkg/handlers/cart`for cart package tests.
#####

### Postman Documentation
- https://documenter.getpostman.com/view/16538634/2s9YJgULVk

### Additional Info About Project (Turkish)
- Projeyi Go dilinde yazdım ve **Gin framework, GORM ve PostgreSQL** kullandım
####
- Package ayrımı olarak, toplamda iki adet olmak üzere **item** ve **cart** packagelerini ayırdım. Aslında ilk önce vas_item ve item'i ayırmayı denedim ancak birbirlerine çok dependent olduklarından dolayı tekrar aynı packageye koydum.
####
- Database tasarımında bir iteme birden çok vas-item eklenebildiğinden ve bir vas-item'ın birden çok itemde bulunabileceğinden dolayı **many-to-many ilişki** ile iki ayrı tablo ve bunları bağlayan bir pivot tablo ile tasarladım.
####
- Database tablolarını yaratırken proje devamlılığını gözeterek **migration dosyaları** kullandım. Bu migrationlarda tabloları yaratırken **primary key, foreign key veya unique** gibi constraintler kullanmadım. Bunun sebebi **soft delete** kullandığım zaman key constraintler soft-deleted entry'lerle çakışmasıydı.
####
- **Controller, Manager, Router, Serializer** yapılarını birbirinden ayırarak ve **dependency injection** ile birbirine geçerek kullanmaya çalıştım.
####
- En başta projeyi tasarladıktan sonra ***integration testleri*** yazarak akışın tamamını görüp daha sonra kodlamaya başladım. Böylede **TDD** mantığını kullanmaya çalıştım.
####
- Methodları kısa tutmaya çalıştım, controllerleri ilk önce çalışacak şekilde uzun yazdım ve daha sonra ufak fonksiyonlara bölüp bu fonksiyonları da **helpers.go** içine koyarak controller'lerin uzunlugunuküçültmeyi denedim. Ancak ***Go'daki error handling***'te sürekli olarak
if check ile error kontrolü yaptığımızdan dolayı ister istemez satır sayısına etkisi oldu.
####
- Birden çok tabloya yazma vs. yapılan controllerlerde, **transaction** yapısını kullanarak **race condition** durumu gerçekleşmesini engelledim.
####
