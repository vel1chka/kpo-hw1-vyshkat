package main

import "fmt"

func main() {
	// 1. Создаем инспектора
	inspector := RealServiceCenter{}
	manager := NewFleetManager(inspector)

	// 2. Подготавливаем тестовые сущности
	goodBike := Bicycle{
		InventoryNumber: InventoryNumber{number: 101},
		Transport:       Transport{difficulty: 2, wear: 5, name: "Городской Велосипед"},
	}

	brokenBike := Bicycle{
		InventoryNumber: InventoryNumber{number: 102},
		Transport:       Transport{difficulty: 3, wear: 35, name: "Сломанный Велосипед"},
	}

	scooter := Scooter{
		InventoryNumber: InventoryNumber{number: 201},
		ElectroTransport: ElectroTransport{
			Transport: Transport{difficulty: 3, wear: 10, name: "Электросамокат Lite"},
			input:     250,
			engine:    "350W",
		},
	}

	fastEBike := EBike{
		InventoryNumber: InventoryNumber{number: 202},
		ElectroTransport: ElectroTransport{
			Transport: Transport{difficulty: 7, wear: 8, name: "Скоростной Электробайк"},
			input:     500,
			engine:    "750W",
		},
	}

	helmet := Helmet{
		InventoryNumber: InventoryNumber{number: 301},
	}

	// 3. Выполняем приемку техники и инвентаря
	fmt.Println("--- Приемка техники ---")

	transports := []TransportableObject{goodBike, brokenBike, scooter, fastEBike}
	for _, t := range transports {
		if err := manager.AddTransport(t); err != nil {
			fmt.Printf("Транспорт №%d (%s) отклонен: %v\n", t.ShowNumber(), t.ShowName(), err)
		} else {
			fmt.Printf("Транспорт №%d (%s) успешно принят в парк\n", t.ShowNumber(), t.ShowName())
		}
	}
	manager.AddItem(helmet)

	// 4. Формируем отчет
	fmt.Println("\n--- Все объекты на балансе ---")
	manager.ShowActivePark()

	fmt.Println("\n--- Энергетический баланс ---")
	fmt.Printf("Суммарный суточный расход: %d кВт·ч\n", manager.TotalInput())

	fmt.Println("\n--- Техника для новичков ---")
	for _, t := range manager.GetBeginnerFriendly() {
		fmt.Printf("Номер: %d | Название: %-22s | Сложность: %d/10\n",
			t.ShowNumber(), t.ShowName(), t.ShowDifficulty())
	}
}
