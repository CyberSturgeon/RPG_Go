package main

import (
	"fmt"
	"io"
)

func checkPointer[T any](p *T) bool {
	return p == nil
}

// TODO: Реализуйте структуры:
type Armor struct {
	Name    string
	Defense int
	Weight  float64
}

func (a Armor) GetName() string {
	return a.Name
}

func (a Armor) GetWeight() float64 {
	return a.Weight
}

func (a *Armor) Use() string {
	if checkPointer(a) {
		return "Предмета нет"
	}
	if a.Defense <= 0 {
		return fmt.Sprintf("Нельзя экипировать %v, предмет сломан", a.GetName())
	}
	a.Defense--
	return fmt.Sprintf("Экипировано %v (%v)", a.GetName(), a.Defense)
}

type Potion struct {
	Name    string
	Effect  string
	Charges int
}

func (p Potion) GetName() string {
	return p.Name
}

func (p Potion) GetWeight() float64 {
	return float64(p.Charges) * 0.5
}

func (p *Potion) Use() string {
	if checkPointer(p) {
		return "Предмета нет"
	}
	if p.Charges <= 0 {
		return fmt.Sprintf("Нельзя выпить %v, предмет выпит", p.GetName())
	}
	p.Charges--
	return fmt.Sprintf("Выпито %v (%v), осталось глотков %v", p.GetName(), p.Effect, p.Charges)
}

type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

func (w Weapon) GetName() string {
	return w.Name
}

func (w Weapon) GetWeight() float64 {
	return float64(w.Damage) * 0.5
}

func (w *Weapon) Use() string {
	if checkPointer(w) {
		return "Предмета нет"
	}
	if w.Durability <= 0 {
		return fmt.Sprintf("Нельзя экипировать %v, предмет сломан", w.GetName())
	}
	w.Durability--
	return fmt.Sprintf("Экипировано %v (%v), прочность %v", w.GetName(), w.Damage, w.Durability)
}

// TODO: Для каждой структуры реализуйте:
// TODO: - Метод Use() string (описание использования, например "Используется <имя>", и изменение Durability или Charges и т.д.)
// TODO: - Методы интерфейса Item

type Item interface {
	GetName() string
	GetWeight() float64
	Use() string
}

// TODO: Реализуйте функцию
func DescribeItem(i Item) string {
	if i == nil {
		return "Предмета нет"
	}
	return fmt.Sprintf("%v (вес: %v)", i.GetName(), i.GetWeight())
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	filtered := make([]T, 0, len(items))
	for _, v := range items {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	out := make([]R, 0, len(items))
	for _, v := range items {
		out = append(out, transform(v))
	}
	return out
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	var zero T
	for _, v := range items {
		if condition(v) {
			return v, true
		}
	}
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	if checkPointer(inv) {
		return
	}
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) GetWeapons() []*Weapon {
	if checkPointer(inv) {
		return nil
	}
	onlyWeapons := Filter(inv.Items, func(it Item) bool {
		_, ok := it.(*Weapon)
		return ok
	})
	return Map(onlyWeapons, func(it Item) *Weapon {
		w, _ := it.(*Weapon)
		return w
	})
}

func (inv *Inventory) GetBrokenItems() []Item {
	if checkPointer(inv) {
		return nil
	}
	return Filter(inv.Items, func(it Item) bool {
		switch v := it.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		case *Armor:
			return v.Defense <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) GetItemNames() []string {
	if checkPointer(inv) {
		return nil
	}
	return Map(inv.Items, func(it Item) string { return it.GetName() })
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(it Item) bool { return it.GetName() == name })
}

// TODO: Бонус: реализуйте интефейс Storable для Weapon и Armor:
// TODO: - Weapon: формат "Weapon|Name|Damage|Durability"
// TODO: - Armor: формат "Armor|Name|Defense|Weight"

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (inv *Inventory) Save(w io.Writer) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func (inv *Inventory) Load(r io.Reader) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
}

func main() {
	sword := &Weapon{Name: "Меч", Damage: 10, Durability: 5}
	shield := &Armor{Name: "Щит", Defense: 5, Weight: 4.5}
	heal := &Potion{Name: "Лечебное", Effect: "+50 HP", Charges: 3}
	brokenBow := &Weapon{Name: "Сломанный лук", Damage: 5, Durability: 0}

	inv := &Inventory{}
	inv.AddItem(sword)
	inv.AddItem(shield)
	inv.AddItem(heal)
	inv.AddItem(brokenBow)

	fmt.Println(sword.Use())
	fmt.Println(shield.Use())
	fmt.Println(heal.Use())
	fmt.Println(brokenBow.Use())

	fmt.Println(DescribeItem(sword))
	fmt.Println(DescribeItem(nil))

	fmt.Println("Оружие в инвентаре:")
	for _, w := range inv.GetWeapons() {
		fmt.Printf("- %v (урон %v, прочность %v)\n", w.Name, w.Damage, w.Durability)
	}

	fmt.Println("Сломанные предметы:")
	for _, it := range inv.GetBrokenItems() {
		fmt.Printf("- %v\n", it.GetName())
	}

	fmt.Println("Названия:", inv.GetItemNames())

	if it, ok := inv.FindItemByName("Щит"); ok {
		fmt.Println("Найден:", DescribeItem(it))
	} else {
		fmt.Println("Щит не найден")
	}
}
