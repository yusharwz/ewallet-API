package blogDto

type (
	GetBlogResponse struct {
		Id          string `json:"id,omitempty"`
		Title       string `json:"title,omitempty"`
		Content     string `json:"content,omitempty"`
		UserId      string `json:"user_id,omitempty"`
		PublishedAt string `json:"published_at,omitempty"`
	}

	CreateBlogRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		UserId  string `json:"user_id"`
	}

	EditBlogRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
)
