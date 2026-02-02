package url

import "URL_checker/internal/repo/entities"

type URLHandler interface {
	postTarget(url string, interval int, timeout int)
	getTargets() []entities.URL
	getTarget(id int) entities.URL
	patchTarget(url string, interval int, timeout int, id int)
	deleteTarget(id int)
}
