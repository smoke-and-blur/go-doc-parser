(function (global) {
    class VNode {
        constructor(tag, props = {}, children = []) {
            this.tag = tag;
            this.props = props;
            this.children = children.flat();
            this.el = null;
        }

        node() {
            const el = document.createElement(this.tag);
            this.el = el;

            // props & events
            for (let [k, v] of Object.entries(this.props)) {
                if (k.startsWith("on") && typeof v === "function") {
                    el[k.toLowerCase()] = v;
                    continue;
                }
                if (v === true) {
                    el.setAttribute(k, "");
                    continue;
                }
                if (v !== false && v != null) {
                    el[k] = v;
                }
            }

            // children
            for (let c of this.children) {
                if (c == null) continue;
                el.appendChild(c.node());
            }

            return el;
        }
    }

    class TextVNode extends VNode {
        constructor(text) {
            super(null);
            this.text = String(text);
        }
        node() {
            this.el = document.createTextNode(this.text);
            return this.el;
        }
    }

    function patchChildren(parent, oldChildren, newChildren) {
        const len = Math.max(oldChildren.length, newChildren.length);

        for (let i = 0; i < len; i++) {
            const oldC = oldChildren[i];
            const newC = newChildren[i];

            if (!oldC && newC) {
                parent.appendChild(newC.node());
                continue;
            }

            if (oldC && !newC) {
                parent.removeChild(oldC.el);
                continue;
            }

            if (oldC instanceof TextVNode && newC instanceof TextVNode) {
                if (oldC.text !== newC.text) {
                    oldC.el.nodeValue = newC.text;
                }
                newC.el = oldC.el;
                continue;
            }

            if (oldC instanceof VNode && newC instanceof VNode) {
                patch(oldC, newC, parent);
            }
        }
    }

    function patch(oldVNode, newVNode, parent) {
        // text nodes are handled in patchChildren
        if (oldVNode.tag !== newVNode.tag) {
            parent.replaceChild(newVNode.node(), oldVNode.el);
            return;
        }

        const el = (newVNode.el = oldVNode.el);

        // update props & events
        for (let [k, v] of Object.entries(newVNode.props)) {
            if (k.startsWith("on") && typeof v === "function") {
                if (oldVNode.props[k] !== v) {
                    el[k.toLowerCase()] = v;
                }
                continue;
            }
            if (oldVNode.props[k] !== v) {
                el[k] = v;
            }
        }

        // diff children
        patchChildren(el, oldVNode.children, newVNode.children);
    }

    function render(component, parent) {
        const newVNode = component();

        if (!parent._vnode) {
            parent._vnode = newVNode;
            parent._component = component;
            parent.appendChild(newVNode.node());
            return;
        }

        patch(parent._vnode, newVNode, parent);
        parent._vnode = newVNode;
    }

    function notify(parent = document.getElementById("app")) {
        console.log("notified")
        render(parent._component, parent);
    }

    // tags proxy
    const tags = new Proxy({}, {
        get(_, name) {
            return (...args) => {
                let props = {};
                let children = [];
                for (const a of args) {
                    if (a == null) continue;
                    if (
                        typeof a === "object" &&
                        !(a instanceof VNode) &&
                        !(a instanceof Node) &&
                        !Array.isArray(a)
                    ) {
                        props = { ...props, ...a };
                        continue;
                    }
                    if (typeof a === "string" || typeof a === "number") {
                        children.push(new TextVNode(a));
                        continue;
                    }
                    children.push(a);
                }
                return new VNode(name, props, children);
            };
        }
    });

    // props proxy
    const props = new Proxy({}, {
        get(_, key) {
            return (strings, ...values) => {
                if (strings.length === 1 && strings[0] === "" && values.length === 0) {
                    return { [key]: true };
                }
                const value = String.raw({ raw: strings }, ...values);
                return { [key]: value };
            };
        }
    });

    // events proxy
    const events = new Proxy({}, {
        get(_, key) {
            if (!key.startsWith("on")) throw new Error("Events must start with 'on'");
            return (fn) => ({ [key]: (e)=>{let r=fn(e);notify();return r} });
        }
    });

    global.erg0 = { VNode, TextVNode, tags, props, events, render, notify };
})(typeof window !== "undefined" ? window : globalThis);
