package converter_test

// Passed on local but failed in ci. Problem with config file path

// func Test_UpdateRate_ShouldGetNewRates(t *testing.T) {
// 	cfg, _ := config.NewFromFile("../../data/config.yaml")
// 	client := rate.New(cfg)
// 	model := converter.New(client)

// 	err := model.UpdateRate()
// 	assert.NoError(t, err)

// 	assert.NotNil(t, model.GetRate())
// }

// func Test_Exchange_ShouldAnswerWithCorrectValue(t *testing.T) {
// 	cfg, _ := config.NewFromFile("../../data/config.yaml")
// 	client := rate.New(cfg)
// 	model := converter.New(client)

// 	amount, err := model.Exchange(1.0, converter.USD, converter.RUB)
// 	assert.NoError(t, err)

// 	assert.Greater(t, amount, 1.0)
// }