// node_modules/prismjs/plugins/line-numbers/prism-line-numbers.js
(function() {
  if (typeof Prism === "undefined" || typeof document === "undefined") {
    return;
  }
  var PLUGIN_NAME = "line-numbers";
  var NEW_LINE_EXP = /\n(?!$)/g;
  var config = Prism.plugins.lineNumbers = {
    /**
     * Get node for provided line number
     *
     * @param {Element} element pre element
     * @param {number} number line number
     * @returns {Element|undefined}
     */
    getLine: function(element, number) {
      if (element.tagName !== "PRE" || !element.classList.contains(PLUGIN_NAME)) {
        return;
      }
      var lineNumberRows = element.querySelector(".line-numbers-rows");
      if (!lineNumberRows) {
        return;
      }
      var lineNumberStart = parseInt(element.getAttribute("data-start"), 10) || 1;
      var lineNumberEnd = lineNumberStart + (lineNumberRows.children.length - 1);
      if (number < lineNumberStart) {
        number = lineNumberStart;
      }
      if (number > lineNumberEnd) {
        number = lineNumberEnd;
      }
      var lineIndex = number - lineNumberStart;
      return lineNumberRows.children[lineIndex];
    },
    /**
     * Resizes the line numbers of the given element.
     *
     * This function will not add line numbers. It will only resize existing ones.
     *
     * @param {HTMLElement} element A `<pre>` element with line numbers.
     * @returns {void}
     */
    resize: function(element) {
      resizeElements([element]);
    },
    /**
     * Whether the plugin can assume that the units font sizes and margins are not depended on the size of
     * the current viewport.
     *
     * Setting this to `true` will allow the plugin to do certain optimizations for better performance.
     *
     * Set this to `false` if you use any of the following CSS units: `vh`, `vw`, `vmin`, `vmax`.
     *
     * @type {boolean}
     */
    assumeViewportIndependence: true
  };
  function resizeElements(elements) {
    elements = elements.filter(function(e) {
      var codeStyles = getStyles(e);
      var whiteSpace = codeStyles["white-space"];
      return whiteSpace === "pre-wrap" || whiteSpace === "pre-line";
    });
    if (elements.length == 0) {
      return;
    }
    var infos = elements.map(function(element) {
      var codeElement = element.querySelector("code");
      var lineNumbersWrapper = element.querySelector(".line-numbers-rows");
      if (!codeElement || !lineNumbersWrapper) {
        return void 0;
      }
      var lineNumberSizer = element.querySelector(".line-numbers-sizer");
      var codeLines = codeElement.textContent.split(NEW_LINE_EXP);
      if (!lineNumberSizer) {
        lineNumberSizer = document.createElement("span");
        lineNumberSizer.className = "line-numbers-sizer";
        codeElement.appendChild(lineNumberSizer);
      }
      lineNumberSizer.innerHTML = "0";
      lineNumberSizer.style.display = "block";
      var oneLinerHeight = lineNumberSizer.getBoundingClientRect().height;
      lineNumberSizer.innerHTML = "";
      return {
        element,
        lines: codeLines,
        lineHeights: [],
        oneLinerHeight,
        sizer: lineNumberSizer
      };
    }).filter(Boolean);
    infos.forEach(function(info) {
      var lineNumberSizer = info.sizer;
      var lines = info.lines;
      var lineHeights = info.lineHeights;
      var oneLinerHeight = info.oneLinerHeight;
      lineHeights[lines.length - 1] = void 0;
      lines.forEach(function(line, index) {
        if (line && line.length > 1) {
          var e = lineNumberSizer.appendChild(document.createElement("span"));
          e.style.display = "block";
          e.textContent = line;
        } else {
          lineHeights[index] = oneLinerHeight;
        }
      });
    });
    infos.forEach(function(info) {
      var lineNumberSizer = info.sizer;
      var lineHeights = info.lineHeights;
      var childIndex = 0;
      for (var i = 0; i < lineHeights.length; i++) {
        if (lineHeights[i] === void 0) {
          lineHeights[i] = lineNumberSizer.children[childIndex++].getBoundingClientRect().height;
        }
      }
    });
    infos.forEach(function(info) {
      var lineNumberSizer = info.sizer;
      var wrapper = info.element.querySelector(".line-numbers-rows");
      lineNumberSizer.style.display = "none";
      lineNumberSizer.innerHTML = "";
      info.lineHeights.forEach(function(height, lineNumber) {
        wrapper.children[lineNumber].style.height = height + "px";
      });
    });
  }
  function getStyles(element) {
    if (!element) {
      return null;
    }
    return window.getComputedStyle ? getComputedStyle(element) : element.currentStyle || null;
  }
  var lastWidth = void 0;
  window.addEventListener("resize", function() {
    if (config.assumeViewportIndependence && lastWidth === window.innerWidth) {
      return;
    }
    lastWidth = window.innerWidth;
    resizeElements(Array.prototype.slice.call(document.querySelectorAll("pre." + PLUGIN_NAME)));
  });
  Prism.hooks.add("complete", function(env) {
    if (!env.code) {
      return;
    }
    var code = (
      /** @type {Element} */
      env.element
    );
    var pre = (
      /** @type {HTMLElement} */
      code.parentNode
    );
    if (!pre || !/pre/i.test(pre.nodeName)) {
      return;
    }
    if (code.querySelector(".line-numbers-rows")) {
      return;
    }
    if (!Prism.util.isActive(code, PLUGIN_NAME)) {
      return;
    }
    code.classList.remove(PLUGIN_NAME);
    pre.classList.add(PLUGIN_NAME);
    var match = env.code.match(NEW_LINE_EXP);
    var linesNum = match ? match.length + 1 : 1;
    var lineNumbersWrapper;
    var lines = new Array(linesNum + 1).join("<span></span>");
    lineNumbersWrapper = document.createElement("span");
    lineNumbersWrapper.setAttribute("aria-hidden", "true");
    lineNumbersWrapper.className = "line-numbers-rows";
    lineNumbersWrapper.innerHTML = lines;
    if (pre.hasAttribute("data-start")) {
      pre.style.counterReset = "linenumber " + (parseInt(pre.getAttribute("data-start"), 10) - 1);
    }
    env.element.appendChild(lineNumbersWrapper);
    resizeElements([pre]);
    Prism.hooks.run("line-numbers", env);
  });
  Prism.hooks.add("line-numbers", function(env) {
    env.plugins = env.plugins || {};
    env.plugins.lineNumbers = true;
  });
})();
//# sourceMappingURL=prismjs_plugins_line-numbers_prism-line-numbers__js.js.map
