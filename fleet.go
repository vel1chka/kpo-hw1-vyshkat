package main

import (
	"errors"
	"fmt"
)

type Numbered interface {
	ShowNumber() int
}

type InventoryNumber struct {
	number int
}

func (i InventoryNumber) ShowNumber() int {
	return i.number
}

type Transportable interface {
	ShowWear() int
	ShowDifficulty() int
	ShowName() string
}

type Transport struct {
	difficulty int
	wear       int
	name       string
}

func (t Transport) ShowWear() int {
	return t.wear
}
func (t Transport) ShowDifficulty() int {
	return t.difficulty
}
func (t Transport) ShowName() string {
	return t.name
}

type UsedElectro interface {
	GetInput() int
}

type ElectroTransport struct {
	Transport
	input  int
	engine string
}

func (e ElectroTransport) GetInput() int {
	return e.input
}

type Bicycle struct {
	InventoryNumber
	Transport
}
type EBike struct {
	InventoryNumber
	ElectroTransport
}
type Scooter struct {
	InventoryNumber
	ElectroTransport
}
type Helmet struct {
	InventoryNumber
}
type Battery struct {
	InventoryNumber
}

type RealServiceCenter struct{}

func (r RealServiceCenter) Check(t Transportable) bool {
	return t.ShowWear() < 20
}

type Inspector interface {
	Check(t Transportable) bool
}

type TransportableObject interface {
	Numbered
	Transportable
}

type FleetManager struct {
	Inspector   Inspector
	ActivePark  []Numbered
	Inventory   []Numbered
	ElectroPark []UsedElectro
	TransPark   []TransportableObject
}

func NewFleetManager(inspector Inspector) *FleetManager {
	return &FleetManager{
		Inspector: inspector,
	}
}
func (f *FleetManager) AddTransport(t TransportableObject) error {
	if !f.Inspector.Check(t) {
		return errors.New("транспорт поломан")
	}

	f.ActivePark = append(f.ActivePark, t)
	f.TransPark = append(f.TransPark, t)

	if electro, ok := t.(UsedElectro); ok {
		f.ElectroPark = append(f.ElectroPark, electro)
	}
	return nil
}
func (f *FleetManager) AddItem(item Numbered) {
	f.ActivePark = append(f.ActivePark, item)
	f.Inventory = append(f.Inventory, item)
}
func (f *FleetManager) ShowActivePark() {
	for _, t := range f.ActivePark {
		fmt.Printf("Номер: %d\n", t.ShowNumber())
	}
}
func (f *FleetManager) TotalInput() int {
	sum := 0
	for _, electro := range f.ElectroPark {
		sum += electro.GetInput()
	}
	return sum
}
func (f *FleetManager) GetBeginnerFriendly() []TransportableObject {
	var result []TransportableObject
	for _, tp := range f.TransPark {
		if tp.ShowDifficulty() <= 4 {
			result = append(result, tp)
		}
	}
	return result
}
