# Overview

Kopycat's GUI is built using HTML, CSS, and JavaScript without any frameworks but using [HTMX](https://htmx.org/) and [jQuery](https://jquery.com/) imported from CDN. 

The main hosted file is `webGUI/dashboard.html` and serves as the entry point for the web application. This file imports all the necessary assests.

## Static Assets

Any static assets, such as images, JavaScript files, or CSS files, that are required by the `webGUI/dashboard.html` file should be placed in the `webGUI/static/` directory. This directory is used to store assets that are imported in the HTML file.

These can be then imported in `webGUI/dashboard.html` using the `<link>` and `<script>` tags from `/static/index.js`.

## Adding New Components

To add new components to the GUI, you can follow these steps:

1. Create any of your new assets in `webGUI/static/`.
2. Import these assets in `webGUI/dashboard.html` using the `<link>` and `<script>` tags.

Easy as this, now your new assets are ready for use in the GUI.
