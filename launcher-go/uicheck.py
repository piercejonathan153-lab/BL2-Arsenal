"""Loads ui.html in headless Chromium with fake launcher state and fails on any script error.
Run before building: python3 uicheck.py  (needs: pip install playwright)"""
import json, sys
from playwright.sync_api import sync_playwright

STATES = [
    dict(installed="1.0.12", latest="1.1.0", krieg="missing"),
    dict(installed="1.1.0", latest="1.0.12", krieg="ready"),
    dict(installed="", latest="1.1.0", krieg="standalone"),
]
html = open("ui.html").read().replace("{{BG}}", "").replace("{{LOGO}}", "")
errors = []
with sync_playwright() as p:
    b = p.chromium.launch()
    for st in STATES:
        s = dict(launcherVersion="0", tps="C:\\TPS", tpsOk=True, sdkOk=True, bl2="C:\\BL2", checked=True,
                 busy=False, launcherUpdate="", notes=["a"], **st)
        mock = "<script>" + "".join(f"window.{n}=()=>Promise.resolve();" for n in
            ["checkUpdates", "installLatest", "installFile", "chooseGame", "openModFolder", "uninstall",
             "updateLauncher", "play", "openUrl", "ready"]) + f"window.getState=()=>Promise.resolve({json.dumps(s)});</script>"
        pg = b.new_page()
        pg.on("pageerror", lambda e: errors.append(str(e)))
        pg.set_content(html.replace("<script>", mock + "<script>", 1))
        pg.wait_for_timeout(300)
        if pg.inner_text("#lv") != "Launcher v0":
            errors.append(f"page did not render for {st}")
    b.close()
print("\n".join(errors) or "ui ok")
sys.exit(1 if errors else 0)
