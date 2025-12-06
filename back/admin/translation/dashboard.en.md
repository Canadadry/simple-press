
## How to Create a Minimal Theme with Markdown Support

In this tutorial, we will create a **minimal theme** using a single base file and a layout file. The article content will be written in **Markdown** and rendered as HTML using the `markdownify` function.

### Step 1: Create the Base HTML File (`baseof.html`)

The `baseof.html` file will contain the overall HTML structure of your page. It will include a placeholder for the content (using `{{template "body" .}}`), which will be populated by the layout file.

Go to the **[Templates](/admin/templates)** section of the admin panel and create a new template called `baseof.html`. This file defines the base structure of your page.

Example content for `baseof.html`:

```html
<html>
  <head></head>
  <body>
    {{template "body" .}}  <!-- This is where the content will be injected -->
  </body>
</html>
```

### Step 2: Create the Layout File (`articles`)

Now create a second template called `articles` in the **[Templates](#)** section. This file will define the content of the page, including the article title and body.

To handle Markdown content, we will use the `markdownify` function inside the layout to convert the Markdown into HTML. The `markdownify` function will be used to process the `Content` field, which is assumed to be written in Markdown.

Example content for `articles`:

```go
{{define "body"}}
  <h1>{{.Title}}</h1>
  <div>{{markdownify .Content}}</div>  <!-- Convert Markdown content to HTML -->
{{end}}
```

In this layout file:

* `{{.Title}}` displays the article title.
* `{{markdownify .Content}}` converts the `Content` field from Markdown to HTML and then renders it inside a `<div>` tag.

### Step 3: Link the Layout to an Article

Now, you can link the layout to an article in your application:

1. Go to the **[Articles](#)** section in the admin panel.
2. Create a new article and fill in the title and content in **Markdown**.
3. Select the layout you created earlier (`articles`).
4. Save the article, and it will automatically be rendered using the layout from `baseof.html` and `articles`.

---

## How to Write Your First Article

Once you’ve created the theme, you can write your first article in Markdown. Follow these steps to create and publish the article:

### Step 1: Access the Articles Section

Click on **[Articles](#)** in the main menu to open the articles management page. This is where you can view and manage all your articles.

### Step 2: Create a New Article

Click on **"Add Article"** in the **[Articles](#)** section. You will need to fill in:

* **Title**: Enter the title of your article, e.g., "My First Article."

* **Content**: Write the content of your article in Markdown. For example:

  ```markdown
  # My First Article

  This is a **simple** article written in *Markdown*.

  ## Subheading

  - Item 1
  - Item 2
  ```

* **Select the Layout**: Choose the layout you created earlier (`main_layout`).

### Step 3: Save the Article

Once you've written the article, click **Save**. The article will be saved and linked to the `main_layout` template, which includes the `markdownify` function to render the Markdown content as HTML.

### Step 4: View the Article

After saving the article, click on its title in the **[Articles](#)** section to view it. The content will be automatically converted from Markdown to HTML and displayed inside the layout defined in `baseof.html`.
