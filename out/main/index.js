"use strict";
const electron = require("electron");
const path = require("path");
let tray = null;
let win = null;
const iconPath = path.join(__dirname, "../../resources/icon.png");
const createWindow = () => {
  win = new electron.BrowserWindow({
    width: 1200,
    height: 800,
    autoHideMenuBar: true,
    webPreferences: {
      preload: path.join(__dirname, "../preload/index.js")
    },
    icon: iconPath
  });
  if (process.env.NODE_ENV === "development") {
    win.loadURL("http://localhost:5173");
  } else {
    win.loadFile(path.join(__dirname, "../../out/renderer/index.html"));
  }
  win.on("close", (event) => {
    event.preventDefault();
    win.hide();
  });
};
function createTray() {
  tray = new electron.Tray(iconPath);
  const contextMenu = electron.Menu.buildFromTemplate([
    {
      label: "退出",
      click: () => {
        electron.app.quit();
        electron.app.exit(0);
      }
    },
    {
      label: "重启",
      click: () => {
        electron.app.relaunch();
        electron.app.exit(0);
      }
    }
  ]);
  tray.setContextMenu(contextMenu);
  tray.setToolTip("WaitingToDo");
  tray.on("click", () => {
    if (win.isVisible()) {
      win.hide();
    } else {
      win.show();
    }
  });
}
electron.app.whenReady().then(() => {
  createWindow();
  createTray();
  electron.app.on("activate", () => {
    if (electron.BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});
electron.app.on("window-all-closed", () => {
  if (process.platform !== "darwin") electron.app.quit();
});
