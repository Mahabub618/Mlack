package models

type BitbucketPayload struct {
	Push struct {
		Changes []struct {
			Commits []struct {
				Hash    string `json:"hash"`
				Message string `json:"message"`
				Author  struct {
					Raw  string `json:"raw"`
					User struct {
						DisplayName string `json:"display_name"`
					} `json:"user"`
				} `json:"author"`
			} `json:"commits"`
		} `json:"changes"`
	} `json:"push"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Actor struct {
		DisplayName string `json:"display_name"`
	} `json:"actor"`
}

type CommitEvent struct {
	Repository string   `json:"repository"`
	Messages   []string `json:"messages"`
	Committer  string   `json:"committer"`
}
