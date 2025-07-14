import { app, BrowserWindow,Tray, Menu, } from 'electron'
import path from 'path'

let tray = null
let win = null
const iconPath = path.join(__dirname, '../../resources/icon.png');
const createWindow = () => {
  win = new BrowserWindow({
    width: 1200,
    height: 800,
    autoHideMenuBar: true,
    webPreferences: {
      preload: path.join(__dirname, '../preload/index.js')
    },
    icon: iconPath
  })

  // 开发环境加载 Vite 开发服务器
  if (process.env.NODE_ENV === 'development') {
    win.loadURL('http://localhost:5173')
    //win.webContents.openDevTools()
  } else {
    // 生产环境加载打包后的文件
    win.loadFile(path.join(__dirname, '../../out/renderer/index.html'))
  }

    // 防止关闭窗口时退出应用
  win.on('close', (event) => {
    event.preventDefault(); // 阻止关闭事件
    win.hide(); // 隐藏窗口
  });
}

// 创建托盘图标
function createTray() {
  tray = new Tray(iconPath); // 托盘图标的路径

  const contextMenu = Menu.buildFromTemplate([
    {
      label: '退出',
      click: () => {
        app.quit(); // 点击退出时退出应用
        app.exit(0);
      },
    },
    {
      label: '重启',
      click: ()=>{
        app.relaunch();
        app.exit(0);
      }
    },
  ]);

  tray.setContextMenu(contextMenu); // 设置托盘菜单
  tray.setToolTip('WaitingToDo'); // 设置悬浮提示

  // 托盘图标点击事件：最小化窗口而不退出
  tray.on('click', () => {
    if (win.isVisible()) {
      win.hide(); // 如果窗口可见，隐藏它
    } else {
      win.show(); // 如果窗口不可见，显示它
    }
  });
}

app.whenReady().then(() => {
  createWindow();
  createTray();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit()
})