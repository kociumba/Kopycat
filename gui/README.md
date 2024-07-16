# Overview

Kopycat's GUI is built using HTML, CSS, and JavaScript without any frameworks but using [HTMX](https://htmx.org/) and [jQuery](https://jquery.com/) imported from CDN. 

The main hosted file is `webGUI/index.html` and serves as the entry point for the web application. This file imports all the necessary assests.

## Static Assets

Any static assets, such as images, JavaScript files, or CSS files, that are required by the `webGUI/index.html` file should be placed in the `webGUI/static/` directory. This directory is used to store assets that are imported in the HTML file.

These can be then imported in `webGUI/index.html` using the `<link>` and `<script>` tags from `/static/index.js`.

## Developing the GUI

There is a vite config for hosting a dev version of the GUI. You can run it locally by installing vite globally with `npm install -g vite` and then running `vite` in the `webGUI` directory.

Or running `vite gui/webGUI` in the `Kopycat` directory which is the same as running `task run-gui` if you have task installed.

### **Keep in mind this dev version will bug out as there is not gonna be a server backing it up, but is great for testing JS functionality and CSS.**

## Adding New Components

To add new components to the GUI, you can follow these steps:

1. Create any of your new assets in `webGUI/static/`.
2. Import these assets in `webGUI/index.html` using the `<link>` and `<script>` tags.

Easy as this, now your new assets are ready for use in the GUI.
