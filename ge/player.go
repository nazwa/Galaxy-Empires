package ge

import (
	"golang.org/x/crypto/bcrypt"
	"sync"
)

const (
	BCRYPT_COMPLEXITY int = 12
)

type PlayerStruct struct {
	planetMutex sync.Mutex `json:"-"`

	ID       string `binding:"omitempty,number"`
	Name     string `form:"name" binding:"required,min=1,max=60"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=5,max=60" json:"-"`
	Planets  []*PlanetStruct
}

type LoginStruct struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"max=200"`
}

func (p *PlayerStruct) GenerateHash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), BCRYPT_COMPLEXITY)
	return string(hash), err
}
func (p *PlayerStruct) HashPassword() error {
	hashedPassword, err := p.GenerateHash(p.Password)
	if err == nil {
		p.Password = hashedPassword
	}
	return err
}

func (p *PlayerStruct) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
}

func (p *PlayerStruct) AddPlanet(planet *PlanetStruct) {
	if planet == nil {
		return
	}
	if planet.Owner == p {
		return
	}

	p.planetMutex.Lock()
	defer p.planetMutex.Unlock()

	if p.Planets == nil {
		p.Planets = make([]*PlanetStruct, 0)
	}
	planet.Owner = p
	p.Planets = append(p.Planets, planet)
}

func (p *PlayerStruct) GetPlanet(id int64) *PlanetStruct {
	p.planetMutex.Lock()
	defer p.planetMutex.Unlock()

	if p.Planets == nil {
		p.Planets = make([]*PlanetStruct, 0)
	}
	if id > int64(len(p.Planets)) {
		return nil
	}
	return p.Planets[id]
}
