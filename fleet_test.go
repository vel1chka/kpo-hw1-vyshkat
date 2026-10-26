package main

import (
	"testing"
)

type FakeInspector struct {
	ShouldPass bool
}

func (f FakeInspector) Check(t Transportable) bool {
	return f.ShouldPass
}

// 1. Тест: Отбраковка сломанного транспорта
func TestAddTransport_RejectsBrokenTransport(t *testing.T) {
	fake := FakeInspector{ShouldPass: false}
	manager := NewFleetManager(fake)

	bike := Bicycle{
		InventoryNumber: InventoryNumber{number: 101},
		Transport:       Transport{difficulty: 2, wear: 50, name: "CityBike"},
	}

	err := manager.AddTransport(bike)
	if err == nil {
		t.Fatalf("Ожидалась ошибка отбраковки инспектором, но получено nil")
	}

	if len(manager.ActivePark) != 0 || len(manager.TransPark) != 0 {
		t.Errorf("Сломанный транспорт не должен попадать в парк")
	}
}

// 2. Тест: Успешное добавление обычного велосипеда
func TestAddTransport_AcceptsMechanicalTransport(t *testing.T) {
	fake := FakeInspector{ShouldPass: true}
	manager := NewFleetManager(fake)

	bike := Bicycle{
		InventoryNumber: InventoryNumber{number: 102},
		Transport:       Transport{difficulty: 3, wear: 5, name: "SportBike"},
	}

	if err := manager.AddTransport(bike); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}

	if len(manager.ActivePark) != 1 || len(manager.TransPark) != 1 {
		t.Fatalf("Транспорт не попал в парки учета")
	}
	if len(manager.ElectroPark) != 0 {
		t.Errorf("Обычный велосипед ошибочно попал в ElectroPark")
	}
}

// 3. Тест: Успешное добавление электросамоката
func TestAddTransport_AcceptsElectroTransport(t *testing.T) {
	fake := FakeInspector{ShouldPass: true}
	manager := NewFleetManager(fake)

	scooter := Scooter{
		InventoryNumber: InventoryNumber{number: 201},
		ElectroTransport: ElectroTransport{
			Transport: Transport{difficulty: 2, wear: 2, name: "Ninebot Max"},
			input:     350,
			engine:    "Electric 350W",
		},
	}

	if err := manager.AddTransport(scooter); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}

	if len(manager.ElectroPark) != 1 || manager.ElectroPark[0].GetInput() != 350 {
		t.Errorf("Самокат не учтен корректно в ElectroPark")
	}
}

// 4. Тест: Добавление инвентаря/вещей
func TestAddItem_AddsDirectlyToInventory(t *testing.T) {
	fake := FakeInspector{ShouldPass: false}
	manager := NewFleetManager(fake)

	helmet := Helmet{InventoryNumber: InventoryNumber{number: 301}}
	battery := Battery{InventoryNumber: InventoryNumber{number: 302}}

	manager.AddItem(helmet)
	manager.AddItem(battery)

	if len(manager.ActivePark) != 2 || len(manager.Inventory) != 2 {
		t.Errorf("Инвентарь не был добавлен корректно")
	}
	if len(manager.TransPark) != 0 {
		t.Errorf("Инвентарь ошибочно попал в TransPark")
	}
}

// 5. Тест: Проверка логики реального сервисного центра
func TestRealServiceCenter_Check(t *testing.T) {
	service := RealServiceCenter{}

	if !service.Check(Transport{wear: 15}) {
		t.Errorf("Износ 15%% должен проходить техосмотр")
	}
	if service.Check(Transport{wear: 20}) {
		t.Errorf("Износ 20%% не должен проходить техосмотр")
	}
	if service.Check(Transport{wear: 45}) {
		t.Errorf("Износ 45%% не должен проходить техосмотр")
	}
}

// 6. Тест поведения: подсчет суммарного расхода энергии
func TestTotalInput_CalculatesSumCorrectly(t *testing.T) {
	fake := FakeInspector{ShouldPass: true}
	manager := NewFleetManager(fake)

	_ = manager.AddTransport(Bicycle{
		InventoryNumber: InventoryNumber{number: 101},
		Transport:       Transport{difficulty: 2, wear: 5, name: "Bike"},
	})
	_ = manager.AddTransport(Scooter{
		InventoryNumber:  InventoryNumber{number: 201},
		ElectroTransport: ElectroTransport{Transport: Transport{wear: 5}, input: 250},
	})
	_ = manager.AddTransport(EBike{
		InventoryNumber:  InventoryNumber{number: 202},
		ElectroTransport: ElectroTransport{Transport: Transport{wear: 5}, input: 500},
	})

	if actual := manager.TotalInput(); actual != 750 {
		t.Errorf("Ожидалось 750, получено %d", actual)
	}
}

// 7. Тест поведения: фильтрация транспорта для новичков
func TestGetBeginnerFriendly_FiltersByDifficulty(t *testing.T) {
	fake := FakeInspector{ShouldPass: true}
	manager := NewFleetManager(fake)

	_ = manager.AddTransport(Bicycle{
		InventoryNumber: InventoryNumber{number: 101},
		Transport:       Transport{difficulty: 2, wear: 5, name: "EasyBike"},
	})
	_ = manager.AddTransport(Scooter{
		InventoryNumber:  InventoryNumber{number: 201},
		ElectroTransport: ElectroTransport{Transport: Transport{difficulty: 4, wear: 5, name: "BorderScooter"}},
	})
	_ = manager.AddTransport(EBike{
		InventoryNumber:  InventoryNumber{number: 202},
		ElectroTransport: ElectroTransport{Transport: Transport{difficulty: 8, wear: 5, name: "ProEBike"}},
	})

	beginners := manager.GetBeginnerFriendly()
	if len(beginners) != 2 {
		t.Fatalf("Ожидалось 2 транспорта для новичков, получено: %d", len(beginners))
	}
	if beginners[0].ShowNumber() != 101 || beginners[0].ShowName() != "EasyBike" {
		t.Errorf("Неверные данные первого элемента выборки")
	}
}

// 8. Тест: Вывод активного парка
func TestShowActivePark_ExecutesWithoutPanic(t *testing.T) {
	fake := FakeInspector{ShouldPass: true}
	manager := NewFleetManager(fake)

	manager.AddItem(Helmet{InventoryNumber: InventoryNumber{number: 301}})
	manager.ShowActivePark()
}
