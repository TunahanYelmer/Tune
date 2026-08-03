package providers


type Track struct{
	Title string
	Artist string
	Album string
	Duration int

}

type Provider interface {
	Login() error
	Play(query string) error
	Pause() error
	Next() error
	Previous() error
	Current() (*Track ,error)
	Search(query string) ([]Track,error)

	SetVolume(level int) error




}
