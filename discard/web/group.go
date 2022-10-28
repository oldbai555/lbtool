package web

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// createStaticHandler create static handler
//func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
//	absolutePath := path.Join(group.prefix, relativePath)
//	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
//	return func(c *Context) error {
//		file := c.Param("filepath")
//		// Check if file exists and/or if we have permission to access it
//		if _, err := fs.Open(file); err != nil {
//			c.Status(http.StatusNotFound)
//			return err
//		}
//
//		fileServer.ServeHTTP(c.Writer, c.Req)
//		return nil
//	}
//}

// Static serve static files
//func (group *RouterGroup) Static(relativePath string, root string) {
//	handler := group.createStaticHandler(relativePath, http.Dir(root))
//	urlPattern := path.Join(relativePath, "/*filepath")
//	// Register GET handlers
//	group.GET(urlPattern, handler)
//}
