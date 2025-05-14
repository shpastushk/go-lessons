package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(int(person.mana) + mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel[0] = byte(int(person.rsel[0]) + respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel[0] = byte(int(person.rsel[0]) + strength<<4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel[1] = byte(int(person.rsel[1]) + experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rsel[1] = byte(int(person.rsel[1]) + level<<4)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(person.mana + 0b0000010000000000)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(person.mana + 0b0001000000000000)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(person.mana + 0b0000100000000000)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(person.mana + uint16(personType)<<13)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name    [42]byte
	mana    uint16 // mana + hasHouse + hasGun + hasFamilty + gType
	health  uint16
	rsel    [2]byte // respect + strength + experience + level
	gold    uint32
	x, y, z int32
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{}
	for _, option := range options {
		option(&gp)
	}
	return gp
}

func (p *GamePerson) Name() string {
	return unsafe.String(unsafe.SliceData(p.name[:]), len(p.name))
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana & 0b0000001111111111)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int(p.rsel[0] & 0b00001111)
}

func (p *GamePerson) Strength() int {
	return int(p.rsel[0] >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.rsel[1] & 0b00001111)
}

func (p *GamePerson) Level() int {
	return int(p.rsel[1] >> 4)
}

func (p *GamePerson) HasHouse() bool {
	return (p.mana & 0b0000010000000000) == 0b0000010000000000
}

func (p *GamePerson) HasGun() bool {
	return (p.mana & 0b0001000000000000) == 0b0001000000000000
}

func (p *GamePerson) HasFamilty() bool {
	return (p.mana & 0b0000100000000000) == 0b0000100000000000
}

func (p *GamePerson) Type() int {
	return int(p.mana & 0b0110000000000000 >> 13)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
