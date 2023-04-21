Dropzone.autoDiscover = false;

const POSTMESSAGE_PREFIX = "filebrowser_";

const getEventName = (name) => (`${POSTMESSAGE_PREFIX}${name}`);
function getExtention(fname) {
  return fname.slice((fname.lastIndexOf(".") - 1 >>> 0) + 2);
}

function pathJoin(parts, sep) {
  var separator = sep || '/';
  var replace = new RegExp(separator + '{1,}', 'g');
  return parts.join(separator).replace(replace, separator);
}

function getQueryString(name) {
  var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
  var r = decodeURIComponent(window.location.search).substr(1).match(reg);
  if (r != null) return r[2].replace(/\+/g, ' ');
  return null;
}

function checkPathNameLegal(name) {
  var reg = new RegExp("[\\/:*<>|]");
  var r = name.match(reg)
  return r == null;
}

function showErrorMessage(jqXHR) {
  let errMsg = jqXHR.getResponseHeader("x-auth-authentication-message")
  if (errMsg == null) {
    errMsg = jqXHR.responseText
  }
  alert(String(jqXHR.status).concat(":", errMsg));
  console.error(errMsg)
}

Vue.filter('fromNow', function (value) {
  return moment(value).fromNow();
})

Vue.filter('formatBytes', function (value) {
  var bytes = parseFloat(value);
  if (bytes < 0) return "-";
  else if (bytes < 1024) return bytes + " B";
  else if (bytes < 1048576) return (bytes / 1024).toFixed(0) + " KB";
  else if (bytes < 1073741824) return (bytes / 1048576).toFixed(1) + " MB";
  else return (bytes / 1073741824).toFixed(1) + " GB";
})

var vm = new Vue({
  el: "#app",
  data: {
    user: {
      email: "",
      name: "",
    },
    homepage: "/",
    location: window.location,
    breadcrumb: [],
    showHidden: false,
    previewMode: false,
    preview: {
      filename: '',
      filetype: '',
      filesize: 0,
      contentHTML: '',
    },
    version: "loading",
    mtimeTypeFromNow: false, // or fromNow
    auth: {},
    search: getQueryString("search"),
    files: [{
      name: "loading ...",
      path: "",
      size: "...",
      type: "dir",
    }],
    myDropzone: null,
  },
  computed: {
    computedFiles: function () {
      var that = this;
      that.preview.filename = null;

      var files = this.files.filter(function (f) {
        if (!that.showHidden && f.name.slice(0, 1) === '.') {
          return false;
        }
        return true;
      });
      return files;
    },
  },
  created: function () {
    $.ajax({
      url: "/-/user",
      method: "get",
      dataType: "json",
      success: function (ret) {
        if (ret) {
          this.user.email = ret.email;
          this.user.name = ret.name;
        }
      }.bind(this)
    })

    const that = this;
    setTimeout(function() {
      that.myDropzone = new Dropzone("#upload-form", {
        paramName: "file",
        maxFilesize: 102400,
        addRemoveLinks: true,
        headers: {
          "X-Requested-File-Server-Token": window.token,
        },
        init: function () {
          this.on("complete", function (file) {
            console.log("reload file list")
            loadFileList()
          })
        }
      });
    }, 1000);
  },
  methods: {
    getLocationPathname: function () {
      return decodeURIComponent(location.pathname);
    },
    getEncodePath: function (filepath, search="") {
      return pathJoin([this.getLocationPathname()].concat(filepath.split("/").map(v => encodeURIComponent(v)))) + search;
    },
    formatTime: function (timestamp) {
      var m = moment(timestamp);
      if (this.mtimeTypeFromNow) {
        return m.fromNow();
      }
      return m.format('YYYY-MM-DD HH:mm:ss');
    },
    toggleHidden: function () {
      this.showHidden = !this.showHidden;
    },
    removeAllUploads: function () {
      this.myDropzone.removeAllFiles();
    },
    parentDirectory: function (path) {
      return path.replace('\\', '/').split('/').slice(0, -1).join('/')
    },
    searchFiles: function () {
      if (!this.search) {
        return
      }
      const search = location.search ? location.search + `&search=${this.search}` : `?search=${this.search}`
      loadFileList(location.pathname + search);
    },
    changeParentDirectory: function (path) {
      var parentDir = this.parentDirectory(path);
      loadFileOrDir(parentDir);
    },
    doubleClickFile: function (f, e) {
      if (f.type === "dir") {
        return
      }
      parent.postMessage({
        event: getEventName(f.type == "dir" ? "dir_double_selected" : "file_double_selected"),
        data: {
          file: f,
        }
      }, '*');
    },
    genDownloadURL: function (f) {
      var search = location.search;
      var sep = search == "" ? "?" : "&"
      return location.origin + this.getEncodePath(f.name) + location.search + sep + "download=true";
    },
    genFileClass: function (f) {
      if (f.type == "dir") {
        if (f.name == '.git') {
          return 'fa-git-square';
        }
        return "fa-folder-open";
      }
      var ext = getExtention(f.name);
      switch (ext) {
        case "go":
        case "py":
        case "js":
        case "java":
        case "c":
        case "cpp":
        case "h":
          return "fa-file-code-o";
        case "pdf":
          return "fa-file-pdf-o";
        case "zip":
          return "fa-file-zip-o";
        case "mp3":
        case "wav":
          return "fa-file-audio-o";
        case "jpg":
        case "png":
        case "gif":
        case "jpeg":
        case "tiff":
          return "fa-file-picture-o";
        case "ipa":
        case "dmg":
          return "fa-apple";
        case "apk":
          return "fa-android";
        case "exe":
          return "fa-windows";
      }
      return "fa-file-text-o"
    },
    clickFileOrDir: function (f, e) {
      parent.postMessage({
        event: getEventName(f.type == "dir" ? "dir_selected" : "file_selected"),
        data: {
          file: f,
        }
      }, '*');
      var reqPath = this.getEncodePath(f.name)
      f.type == "dir" && loadFileOrDir(reqPath);
      e.preventDefault()
    },
    changePath: function (reqPath, e) {
      reqPath = reqPath.replace(/\/+/g, "/")
      parent.postMessage({
        event: getEventName("path_changed"),
        data: {
          path: reqPath,
        }
      }, '*');
      loadFileOrDir(reqPath);
      e.preventDefault()
    },
    backOnePath: function (e) {
      const currentPath = this.getLocationPathname().replace(/\/+$/, "");
      const homepage = this.homepage.replace(/\/+$/, "");
      if (currentPath === homepage) {
        return
      }
      const reqPath = this.parentDirectory(this.getLocationPathname()).replace(/\/+$/, "");
      if (reqPath.length < homepage.length) {
        return;
      }
      this.changePath(reqPath, e);
    },
    showInfo: function (f) {
      $.ajax({
        url: this.getEncodePath(f.name),
        data: {
          op: "info",
        },
        method: "GET",
        success: function (res) {
          $("#file-info-title").text(f.name);
          $("#file-info-content").text(JSON.stringify(res, null, 4));
          $("#file-info-modal").modal("show");
          // console.log(JSON.stringify(res, null, 4));
        },
        error: function (jqXHR, textStatus, errorThrown) {
          showErrorMessage(jqXHR)
        }
      })
    },
    makeDirectory: function () {
      var name = window.prompt("current path: " + this.getLocationPathname() + "\nplease enter the new directory name", "")
      // console.log(name)
      if (!name) {
        return
      }
      parent.postMessage({
        event: getEventName("dir_created"),
        data: {
          name: name,
        }
      }, '*');
      if(!checkPathNameLegal(name)) {
        alert("Name should not contains any of \\/:*<>|")
        return
      }
      $.ajax({
        url: this.getEncodePath(name) + "?type=folder",
        method: "POST",
        success: function (res) {
          // console.log(res)
          loadFileList()
        },
        error: function (jqXHR, textStatus, errorThrown) {
          showErrorMessage(jqXHR)
        }
      })
    },
    deletePathConfirm: function (f, e) {
      e.preventDefault();
      parent.postMessage({
        event: getEventName("dir_deleted"),
        data: {
          file: f
        }
      }, '*');
      if (!e.altKey) { // skip confirm when alt pressed
        if (!window.confirm("Delete " + f.name + " ?")) {
          return;
        }
      }
      $.ajax({
        url: this.getEncodePath(f.name),
        method: 'DELETE',
        success: function (res) {
          loadFileList()
        },
        error: function (jqXHR, textStatus, errorThrown) {
          showErrorMessage(jqXHR)
        }
      });
    },
    updateBreadcrumb: function (pathname) {
      var pathname = decodeURIComponent(pathname || this.getLocationPathname() || "/");
      pathname = pathname.split('?')[0]
      var parts = pathname.split('/');
      this.breadcrumb = [];
      if (pathname == "/") {
        return this.breadcrumb;
      }
      var i = 2;
      for (; i <= parts.length; i += 1) {
        var name = parts[i - 1];
        if (!name) {
          continue;
        }
        var path = parts.slice(0, i).join('/');
        this.breadcrumb.push({
          name: name + (i == parts.length ? ' /' : ''),
          path: path
        })
      }
      return this.breadcrumb;
    },
    loadPreviewFile: function (filepath, e) {
      if (e) {
        e.preventDefault() // may be need a switch
      }
      var that = this;
      $.getJSON(pathJoin(['/-/info', this.getLocationPathname()]))
          .then(function (res) {
            // console.log(res);
            that.preview.filename = res.name;
            that.preview.filesize = res.size;
            return $.ajax({
              url: '/' + res.path,
              dataType: 'text',
            });
          })
          .then(function (res) {
            // console.log(res)
            that.preview.contentHTML = '<pre>' + res + '</pre>';
            // console.log("Finally")
          })
          .done(function (res) {
            // console.log("done", res)
          });
    },
    loadAll: function () {
      // TODO: move loadFileList here
    },
  }
})

window.onpopstate = function (event) {
  if (location.search.match(/\?search=/)) {
    location.reload();
    return;
  }
  loadFileList()
}

function loadFileOrDir(reqPath) {
  const requestUri = (reqPath || "/") + location.search
  const retObj = loadFileList(requestUri)
  if (retObj !== null) {
    retObj.done(function () {
      window.history.pushState({}, "", requestUri);
    });
  }

}

function loadFileList(pathname) {
  var pathname = pathname || decodeURIComponent(location.pathname) + location.search;
  var retObj = null
  // TODO: rewrite the type of raw
  if (getQueryString("raw") !== "false") { // not a file preview
    var sep = pathname.indexOf("?") === -1 ? "?" : "&"
    retObj = $.ajax({
      url: pathname + sep + "json=true",
      dataType: "json",
      cache: false,
      success: function (res) {
        res.files = _.sortBy(res.files, function (f) {
          var weight = f.type == 'dir' ? 1000 : 1;
          return -weight * f.mtime;
        })
        vm.files = res.files;
        vm.auth = res.auth;
        const configs = res.configs;
        prefixReflect = configs?.prefixReflect;
        pathname = decodeURIComponent(pathname)
        if (prefixReflect && prefixReflect.length > 0) {
          for (let i = 0; i < prefixReflect.length; i++) {
            if (pathname.startsWith(prefixReflect[i])) {
              vm.homepage = prefixReflect[i].startsWith('/') ? prefixReflect[i] : '/' + prefixReflect[i];
              pathname = pathname.replace(prefixReflect[i], '')
              break;
            }
            re = new RegExp(prefixReflect[i])
            res = pathname.match(re)
            if (res) {
              vm.homepage = res[0].startsWith('/') ? res[0] : '/' + res[0]
              pathname = pathname.replace(re, '')
              break;
            }
          }
          if (!pathname.startsWith('/')) {
            pathname = '/' + pathname
          }
        }
        vm.updateBreadcrumb(pathname);
      },
      error: function (jqXHR, textStatus, errorThrown) {
        showErrorMessage(jqXHR)
      },
    });

  }

  // vm.previewMode = getQueryString("raw") == "false";
  // if (vm.previewMode) {
  //   vm.loadPreviewFile();
  // }
  return retObj
}

$(function () {
  $.ajaxSetup({
    beforeSend: function (xhr) {
      window.token && xhr.setRequestHeader("X-Requested-File-Server-Token", window.token);
    }
  })

  $.scrollUp({
    scrollText: '', // text are defined in css
  });

  // For page first loading
  loadFileList(decodeURIComponent(location.pathname) + location.search)

  // update version
  $.getJSON("/-/sysinfo", function (res) {
    vm.version = res.version;
  })

  var clipboard = new Clipboard('.btn');
  clipboard.on('success', function (e) {
    console.info('Action:', e.action);
    console.info('Text:', e.text);
    console.info('Trigger:', e.trigger);
    $(e.trigger)
        .tooltip('show')
        .mouseleave(function () {
          $(this).tooltip('hide');
        })

    e.clearSelection();
  });
});
