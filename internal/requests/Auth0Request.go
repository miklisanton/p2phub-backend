package requests

type Auth0Request struct {
    Email    string `json:"email"`
    Secret string `json:"secret"`
}
