V1
- [x] list file
- [x] list block
- [x] edit block
- [x] edit layout of article
- [x] add bock to article
- [x] edit bock data of article
  - [x] render block of data
  - [x] can edit block of data on the frontend
  - [x] can edit block of data on the backend
  - [x] connect them
  - [x] create a layout that display block
- [x] remove front SSR
- [x] add file, article,  block,template, layout
- [x] fetch article from go template to allow paginated list of articles
- [] edit filename, delete file

V1.1 export
- [] add figure to all article (allow upload file from article, or open file selector)
- [] define front url in .env
- [] serve front url from back
- [] add function to allow fetching page name, url, image, content ???
- [] write a build static site script around wget
- [] deploy a first article on private server
- [] expose build static as route so downloading site from front become possible

V1.1.1 polish
- [] remove file, article,block in article, block,template, layout
- [] list bug and fix them
- [] basic unit test around bug of renaming template or layout in front

V1.2 SEO
- [] duplicate article, layout, template
- [] add seo variable (or dynmaque data, like blog but for an article so only one definition)
- [] allow resizing image with a query url request

V2
- [] markdown editor on article content edition
- [] swap go template engine for otto + personal jsx parser
- [] redesign block data definition => list of fullpath + option (toggle stuff, select other)

V3 light
- [] remove npm and use upkg or else and complie with raw esbuild call
