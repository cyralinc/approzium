const express = require(`express`)

// Enable development support for serving HTML from `./static` folder
// reference: https://github.com/gatsbyjs/gatsby/issues/13072
exports.onCreateDevServer = ({ app }) => {
    app.use(express.static(`public`))
}
