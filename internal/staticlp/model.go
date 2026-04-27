package staticlp

// Campaign is the on-disk shape of data/campaigns/<id>.toml (go-toml round-trip).
type Campaign struct {
	ID                 string     `toml:"id"`
	FormName           string     `toml:"form_name"`
	HoneypotField      string     `toml:"honeypot_field"`
	SubmitLabel        string     `toml:"submit_label"`
	ConsentText        string     `toml:"consent_text"`
	QRDelayMS          int        `toml:"qr_delay_ms"`
	QRMetaRefreshSec   int        `toml:"qr_meta_refresh_sec"`
	ThankYouRedirectMS int        `toml:"thank_you_redirect_ms"`
	Paths              Paths      `toml:"paths"`
	Delivery           Delivery   `toml:"delivery"`
	Presenter          *Presenter `toml:"presenter,omitempty"`
	Analytics          Analytics  `toml:"analytics"`
	Fields             []Field    `toml:"fields"`
}

type Paths struct {
	ThankYou string `toml:"thank_you"`
	Material string `toml:"material"`
	Landing  string `toml:"landing"`
}

type Delivery struct {
	Type       string `toml:"type"`
	FileID     string `toml:"file_id"`
	PreviewURL string `toml:"preview_url"`
	DirectURL  string `toml:"direct_url"`
}

type Presenter struct {
	Name        string `toml:"name"`
	Headline    string `toml:"headline"`
	Bio         string `toml:"bio"`
	Photo       string `toml:"photo"`
	LinkedinURL string `toml:"linkedin_url"`
}

type Analytics struct {
	BodyClass string `toml:"body_class"`
	Campaign  string `toml:"campaign"`
}

type Field struct {
	Name         string `toml:"name"`
	Label        string `toml:"label"`
	Type         string `toml:"type"`
	Required     bool   `toml:"required"`
	Autocomplete string `toml:"autocomplete"`
	Inputmode    string `toml:"inputmode"`
}
