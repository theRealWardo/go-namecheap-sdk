package namecheap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersGetPricing(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="UTF-8"?>
		<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
			<Errors />
			<RequestedCommand>namecheap.users.getPricing</RequestedCommand>
			<CommandResponse Type="namecheap.users.getPricing">
				<UserGetPricingResult>
					<ProductType Name="DOMAIN">
						<ProductCategory Name="REACTIVATE">
							<Product Name="biz">
								<Price Duration="1" DurationType="YEAR" Price="8.55" RegularPrice="8.55" YourPrice="8.55" CouponPrice="" Currency="USD" />
								<Price Duration="2" DurationType="YEAR" Price="8.87" RegularPrice="8.87" YourPrice="8.87" CouponPrice="" Currency="USD" />
							</Product>
							<Product Name="bz">
								<Price Duration="1" DurationType="YEAR" Price="8.88" RegularPrice="8.88" YourPrice="8.88" CouponPrice="" Currency="USD" />
							</Product>
						</ProductCategory>
						<ProductCategory Name="REGISTER">
							<Product Name="biz">
								<Price Duration="1" DurationType="YEAR" Price="6.00" RegularPrice="6.00" YourPrice="6.00" CouponPrice="" Currency="USD" />
								<Price Duration="2" DurationType="YEAR" Price="8.87" RegularPrice="8.87" YourPrice="8.87" CouponPrice="" Currency="USD" />
							</Product>
						</ProductCategory>
					</ProductType>
				</UserGetPricingResult>
			</CommandResponse>
			<Server>IMWS-A06</Server>
			<GMTTimeDifference>+5:30</GMTTimeDifference>
			<ExecutionTime>1.109</ExecutionTime>
		</ApiResponse>
	`

	t.Run("request_command", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		args := &GetPricingArgs{
			ProductType: ProductTypeDomain,
		}

		_, err := client.Users.GetPricing(args)
		if err != nil {
			t.Fatal("Error calling GetPricing", err)
		}

		assert.Equal(t, "namecheap.users.getPricing", sentBody.Get("Command"))
		assert.Equal(t, "DOMAIN", sentBody.Get("ProductType"))
	})

	t.Run("parse_response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		args := &GetPricingArgs{
			ProductType: ProductTypeDomain,
		}

		result, err := client.Users.GetPricing(args)
		if err != nil {
			t.Fatal("Error calling GetPricing", err)
		}

		assert.NotNil(t, result)
		assert.NotNil(t, result.UserGetPricingResult)
		assert.NotNil(t, result.UserGetPricingResult.ProductTypes)
		assert.Equal(t, 1, len(*result.UserGetPricingResult.ProductTypes))

		productType := (*result.UserGetPricingResult.ProductTypes)[0]
		assert.Equal(t, "DOMAIN", *productType.Name)
		assert.NotNil(t, productType.ProductCategories)
		assert.Equal(t, 2, len(*productType.ProductCategories))

		reactivateCategory := (*productType.ProductCategories)[0]
		assert.Equal(t, "REACTIVATE", *reactivateCategory.Name)
		assert.NotNil(t, reactivateCategory.Products)
		assert.Equal(t, 2, len(*reactivateCategory.Products))

		bizProduct := (*reactivateCategory.Products)[0]
		assert.Equal(t, "biz", *bizProduct.Name)
		assert.NotNil(t, bizProduct.Prices)
		assert.Equal(t, 2, len(*bizProduct.Prices))

		price := (*bizProduct.Prices)[0]
		assert.Equal(t, "1", *price.Duration)
		assert.Equal(t, "YEAR", *price.DurationType)
		assert.Equal(t, "8.55", *price.Price)
		assert.Equal(t, "8.55", *price.RegularPrice)
		assert.Equal(t, "8.55", *price.YourPrice)
		assert.Equal(t, "USD", *price.Currency)
	})

	t.Run("with_optional_parameters", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := io.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		productCategory := ProductCategoryDomains
		actionName := ActionNameRegister
		productName := ProductNameCom
		promoCode := "TESTPROMO"

		args := &GetPricingArgs{
			ProductType:     ProductTypeDomain,
			ProductCategory: &productCategory,
			ActionName:      &actionName,
			ProductName:     &productName,
			PromotionCode:   &promoCode,
		}

		_, err := client.Users.GetPricing(args)
		if err != nil {
			t.Fatal("Error calling GetPricing", err)
		}

		assert.Equal(t, "namecheap.users.getPricing", sentBody.Get("Command"))
		assert.Equal(t, "DOMAIN", sentBody.Get("ProductType"))
		assert.Equal(t, "DOMAINS", sentBody.Get("ProductCategory"))
		assert.Equal(t, "REGISTER", sentBody.Get("ActionName"))
		assert.Equal(t, "COM", sentBody.Get("ProductName"))
		assert.Equal(t, "TESTPROMO", sentBody.Get("PromotionCode"))
	})

	t.Run("validation_missing_product_type", func(t *testing.T) {
		client := setupClient(nil)

		args := &GetPricingArgs{}

		_, err := client.Users.GetPricing(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ProductType is required")
	})
}
