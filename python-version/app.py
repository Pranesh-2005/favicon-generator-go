import io
import zipfile
import json
import time
import tempfile
import os
from PIL import Image
import gradio as gr

def center_crop_square(img: Image.Image) -> Image.Image:
    w, h = img.size
    min_side = min(w, h)
    left = (w - min_side) // 2
    top = (h - min_side) // 2
    right = left + min_side
    bottom = top + min_side
    return img.crop((left, top, right, bottom))

def generate_favicons(image, name, short_name, theme_color, background_color, tile_color):
    if image is None:
        return None

    img = Image.open(image).convert("RGBA")
    img = center_crop_square(img)

    # Ensure a very high-res base image
    base_size = max(1024, img.width, img.height)
    img = img.resize((base_size, base_size), Image.LANCZOS)

    if not name:
        name = "My App"
    if not short_name:
        short_name = "App"
    if not theme_color:
        theme_color = "#ffffff"
    if not background_color:
        background_color = "#ffffff"
    if not tile_color:
        tile_color = theme_color

    # Icon specs: regular + ultra-high for Apple/Android
    icon_specs = [
        (16, "favicon-16x16.png"),
        (32, "favicon-32x32.png"),
        (48, "favicon-48x48.png"),
        (57, "apple-icon-57x57.png"),
        (60, "apple-icon-60x60.png"),
        (72, "apple-icon-72x72.png"),
        (76, "apple-icon-76x76.png"),
        (96, "android-icon-96x96.png"),
        (114, "apple-icon-114x114.png"),
        (120, "apple-icon-120x120.png"),
        (128, "android-icon-128x128.png"),
        (144, "android-icon-144x144.png"),
        (152, "apple-icon-152x152.png"),
        (167, "apple-icon-167x167.png"),
        (180, "apple-icon-180x180.png"),
        (192, "android-icon-192x192.png"),
        (256, "android-icon-256x256.png"),
        (384, "android-icon-384x384.png"),
        (512, "android-icon-512x512.png"),
        (1024, "apple-icon-1024x1024.png"),  # Retina / App Store
        (150, "ms-icon-150x150.png"),
        (32, "ms-icon-32x32.png"),
        (70, "ms-icon-70x70.png"),
        (310, "ms-icon-310x310.png"),
    ]

    files = {}

    # Generate PNGs
    for size, filename in icon_specs:
        resized = img.resize((size, size), Image.LANCZOS)
        buf = io.BytesIO()
        resized.save(buf, format="PNG")
        files[filename] = buf.getvalue()

    # Apple touch icon
    files["apple-touch-icon.png"] = files.get("apple-icon-180x180.png", img.resize((180,180), Image.LANCZOS).tobytes())

    # favicon.ico (multi-size)
    ico_sizes = [16, 32, 48, 64, 128, 256]
    ico_images = [img.resize((s, s), Image.LANCZOS) for s in ico_sizes]
    ico_buf = io.BytesIO()
    ico_images[0].save(ico_buf, format='ICO', sizes=[(s, s) for s in ico_sizes])
    files["favicon.ico"] = ico_buf.getvalue()

    # manifest.json
    icons = [
        {
            "src": fn,
            "sizes": f"{sz}x{sz}",
            "type": "image/png"
        }
        for sz, fn in icon_specs if fn.endswith(".png") and sz <= 512
    ]
    icons.append({
        "src": "android-icon-512x512.png",
        "sizes": "512x512",
        "type": "image/png",
        "purpose": "any maskable"
    })

    manifest = {
        "name": name,
        "short_name": short_name,
        "start_url": "/",
        "display": "standalone",
        "background_color": background_color,
        "theme_color": theme_color,
        "icons": icons
    }
    files["site.webmanifest"] = json.dumps(manifest, indent=2).encode()
    files["manifest.json"] = json.dumps(manifest, indent=2).encode()

    # browserconfig.xml
    bc = f"""<?xml version="1.0" encoding="utf-8"?>
<browserconfig><msapplication><tile>
  <square70x70logo src="ms-icon-70x70.png"/>
  <square150x150logo src="ms-icon-150x150.png"/>
  <square310x310logo src="ms-icon-310x310.png"/>
  <TileColor>{tile_color}</TileColor>
</tile></msapplication></browserconfig>
"""
    files["browserconfig.xml"] = bc.encode()

    # README.txt
    readme = f"Generated favicons for {name}\n\nAll icons are in the root of the ZIP.\nTheme color: {theme_color}\n"
    files["README.txt"] = readme.encode()

    # ZIP
    zip_buf = io.BytesIO()
    with zipfile.ZipFile(zip_buf, "w") as zf:
        for path, data in files.items():
            zf.writestr(path, data)
    zip_buf.seek(0)

    ts = int(time.time())
    filename = f"favicons-{ts}.zip"
    return filename, zip_buf

def gradio_ui(image, name, short_name, theme_color, background_color, tile_color):
    result = generate_favicons(image, name, short_name, theme_color, background_color, tile_color)
    if result is None:
        return None
    filename, zip_buf = result
    
    # Save ZIP to temporary file for Gradio download
    tmp_dir = tempfile.mkdtemp()
    filepath = os.path.join(tmp_dir, filename)
    with open(filepath, "wb") as f:
        f.write(zip_buf.read())
    
    return filepath

demo = gr.Interface(
    fn=gradio_ui,
    inputs=[
        gr.Image(type="filepath", label="Upload Image"),
        gr.Textbox(label="App Name", placeholder="My App"),
        gr.Textbox(label="Short Name", placeholder="App"),
        gr.ColorPicker(label="Theme Color", value="#ffffff"),
        gr.ColorPicker(label="Background Color", value="#ffffff"),
        gr.ColorPicker(label="Tile Color", value="#ffffff"),
    ],
    outputs=gr.File(label="Download Favicons ZIP"),
    title="Ultra High-Quality Favicon Generator",
    description="Upload a high-res image (≥1024px) and generate ultra-high quality favicons for Retina and modern devices."
)

if __name__ == "__main__":
    demo.launch(server_name="0.0.0.0", server_port=7860, debug=True, pwa=True)