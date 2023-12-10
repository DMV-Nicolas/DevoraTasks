package util

import (
	"fmt"
	"math/rand"
)

func RandomUsername() string {
	monsters := []string{
		"Caballo", "Omnitorrinco", "Avion", "HijoDePuta", "Avi",
		"Luis", "Rodrigo", "Andres", "Santiago", "Diego", "Gustavo",
		"Juan", "Nicolas", "Cristiansito", "Juliansito", "Valentinita",
		"Sebastian", "David",
	}
	actions := []string{
		"VioladorDe", "AbusadorDe", "TerapeutaDe", "SimpDe", "DesarmadorDe",
		"DominadoPor", "EsclavizadoPor", "SexualmenteAbusadoPor", "AtraidoPor",
		"AmanteDelSexoCon", "AsesinadoPor", "TraumadoPor", "LiderDe", "CreadorDe",
	}
	victims := []string{
		"Abuelas", "Feministas", "Comunistas", "Capitalistas", "Langostas",
		"Hombres", "Jirafas", "Penes", "Duendes",
	}
	str := monsters[rand.Intn(len(monsters))]
	str += actions[rand.Intn(len(actions))]
	str += victims[rand.Intn(len(victims))]
	str += fmt.Sprint(rand.Intn(1000))
	return str
}

func RandomEmail() string {
	names := []string{
		"abuela", "cristian", "santiago",
		"feminista", "nicolas", "juan",
		"comunista", "julian", "diego",
		"capitalista", "pepito", "valentina",
		"langosta", "avi", "maria",
		"hombre", "rodolfo", "fernando",
		"jirafa", "gustavo", "proplayer",
		"pene", "rodrigo", "noob",
		"duende", "luis", "hacker",
	}
	business := []string{
		"gmail", "google", "yt",
		"unal.edu", "colsubsidio.edu", "hotmail",
		"outlook", "sakura",
	}
	countries := []string{
		"com", "co", "ar", "es", "us",
		"br", "cl", "pe", "mx", "uy",
	}
	str := names[rand.Intn(len(names))]
	str += fmt.Sprint(rand.Intn(1000))
	str += "@"
	str += business[rand.Intn(len(business))] + "."
	str += countries[rand.Intn(len(countries))]
	return str
}

func RandomPassword(size int) (str string) {
	digits := "abcdefghijklmnopqrstuvwxyz1234567890"
	for i := 0; i < size; i++ {
		str += string(digits[rand.Intn(len(digits))])
	}
	return str
}

func RandomTitle() string {
	actions := []string{
		"Violar", "Abusar", "Adoptar", "Cocinarle", "Cogerme",
		"Hacerle una mamada", "Matar", "Incendiar", "Darle droga",
	}
	names := []string{
		"mi abuela", "cristian", "santiago",
		"una feminista", "nicolas", "juan",
		"un comunista", "julian", "diego",
		"un capitalista", "pepito", "valentina",
		"una langosta", "avi", "maria",
		"un hombre", "rodolfo", "fernando",
		"una jirafa", "gustavo", "un proplayer",
		"pene", "rodrigo", "un noob",
		"un duende", "luis", "un hacker",
	}
	str := actions[rand.Intn(len(actions))] + " a "
	str += names[rand.Intn(len(names))]
	return str
}

// RandomInt generates a random int number between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomID() int64 {
	return int64(RandomInt(1, 1000))
}
