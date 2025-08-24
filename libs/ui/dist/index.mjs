var __defProp = Object.defineProperty;
var __name = (target, value) => __defProp(target, "name", { value, configurable: true });

// src/components/ui/button.tsx
import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva } from "class-variance-authority";

// src/lib/utils.ts
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";
function cn(...inputs) {
  return twMerge(clsx(inputs));
}
__name(cn, "cn");

// src/components/ui/button.tsx
var buttonVariants = cva("inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50", {
  variants: {
    variant: {
      default: "bg-primary text-primary-foreground hover:bg-primary/90",
      destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
      outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
      secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
      ghost: "hover:bg-accent hover:text-accent-foreground",
      link: "text-primary underline-offset-4 hover:underline"
    },
    size: {
      default: "h-10 px-4 py-2",
      sm: "h-9 rounded-md px-3",
      lg: "h-11 rounded-md px-8",
      icon: "h-10 w-10"
    }
  },
  defaultVariants: {
    variant: "default",
    size: "default"
  }
});
var Button = /* @__PURE__ */ React.forwardRef(({ className, variant, size, asChild = false, ...props }, ref) => {
  const Comp = asChild ? Slot : "button";
  return /* @__PURE__ */ React.createElement(Comp, {
    className: cn(buttonVariants({
      variant,
      size,
      className
    })),
    ref,
    ...props
  });
});
Button.displayName = "Button";

// src/components/ui/card.tsx
import * as React2 from "react";
var Card = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("div", {
  ref,
  className: cn("rounded-lg border bg-card text-card-foreground shadow-sm", className),
  ...props
}));
Card.displayName = "Card";
var CardHeader = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("div", {
  ref,
  className: cn("flex flex-col space-y-1.5 p-6", className),
  ...props
}));
CardHeader.displayName = "CardHeader";
var CardTitle = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("h3", {
  ref,
  className: cn("text-2xl font-semibold leading-none tracking-tight", className),
  ...props
}));
CardTitle.displayName = "CardTitle";
var CardDescription = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("p", {
  ref,
  className: cn("text-sm text-muted-foreground", className),
  ...props
}));
CardDescription.displayName = "CardDescription";
var CardContent = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("div", {
  ref,
  className: cn("p-6 pt-0", className),
  ...props
}));
CardContent.displayName = "CardContent";
var CardFooter = /* @__PURE__ */ React2.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React2.createElement("div", {
  ref,
  className: cn("flex items-center p-6 pt-0", className),
  ...props
}));
CardFooter.displayName = "CardFooter";

// src/components/ui/badge.tsx
import * as React3 from "react";
import { cva as cva2 } from "class-variance-authority";
var badgeVariants = cva2("inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2", {
  variants: {
    variant: {
      default: "border-transparent bg-primary text-primary-foreground hover:bg-primary/80",
      secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
      destructive: "border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80",
      outline: "text-foreground",
      success: "border-transparent bg-green-500 text-white hover:bg-green-500/80",
      warning: "border-transparent bg-yellow-500 text-white hover:bg-yellow-500/80",
      info: "border-transparent bg-blue-500 text-white hover:bg-blue-500/80"
    }
  },
  defaultVariants: {
    variant: "default"
  }
});
function Badge({ className, variant, ...props }) {
  return /* @__PURE__ */ React3.createElement("div", {
    className: cn(badgeVariants({
      variant
    }), className),
    ...props
  });
}
__name(Badge, "Badge");

// src/components/ui/input.tsx
import * as React4 from "react";
var Input = /* @__PURE__ */ React4.forwardRef(({ className, type, ...props }, ref) => {
  return /* @__PURE__ */ React4.createElement("input", {
    type,
    className: cn("flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50", className),
    ref,
    ...props
  });
});
Input.displayName = "Input";

// src/components/ui/label.tsx
import * as React5 from "react";
import * as LabelPrimitive from "@radix-ui/react-label";
import { cva as cva3 } from "class-variance-authority";
var labelVariants = cva3("text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70");
var Label = /* @__PURE__ */ React5.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React5.createElement(LabelPrimitive.Root, {
  ref,
  className: cn(labelVariants(), className),
  ...props
}));
Label.displayName = LabelPrimitive.Root.displayName;

// src/components/ui/toast.tsx
import * as React6 from "react";
import * as ToastPrimitives from "@radix-ui/react-toast";
import { cva as cva4 } from "class-variance-authority";
import { X } from "lucide-react";
var ToastProvider = ToastPrimitives.Provider;
var ToastViewport = /* @__PURE__ */ React6.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React6.createElement(ToastPrimitives.Viewport, {
  ref,
  className: cn("fixed top-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-0 sm:right-0 sm:top-auto sm:flex-col md:max-w-[420px]", className),
  ...props
}));
ToastViewport.displayName = ToastPrimitives.Viewport.displayName;
var toastVariants = cva4("group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=move]:transition-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[swipe=end]:animate-out data-[state=closed]:fade-out-80 data-[state=closed]:slide-out-to-right-full data-[state=open]:slide-in-from-top-full data-[state=open]:sm:slide-in-from-bottom-full", {
  variants: {
    variant: {
      default: "border bg-background text-foreground",
      destructive: "destructive border-destructive bg-destructive text-destructive-foreground",
      success: "border-green-200 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-900 dark:text-green-100"
    }
  },
  defaultVariants: {
    variant: "default"
  }
});
var Toast = /* @__PURE__ */ React6.forwardRef(({ className, variant, ...props }, ref) => {
  return /* @__PURE__ */ React6.createElement(ToastPrimitives.Root, {
    ref,
    className: cn(toastVariants({
      variant
    }), className),
    ...props
  });
});
Toast.displayName = ToastPrimitives.Root.displayName;
var ToastAction = /* @__PURE__ */ React6.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React6.createElement(ToastPrimitives.Action, {
  ref,
  className: cn("inline-flex h-8 shrink-0 items-center justify-center rounded-md border bg-transparent px-3 text-xs font-medium ring-offset-background transition-colors hover:bg-secondary focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 group-[.destructive]:border-muted/40 group-[.destructive]:hover:border-destructive/30 group-[.destructive]:hover:bg-destructive group-[.destructive]:hover:text-destructive-foreground group-[.destructive]:focus:ring-destructive", className),
  ...props
}));
ToastAction.displayName = ToastPrimitives.Action.displayName;
var ToastClose = /* @__PURE__ */ React6.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React6.createElement(ToastPrimitives.Close, {
  ref,
  className: cn("absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100 group-[.destructive]:text-red-300 group-[.destructive]:hover:text-red-50 group-[.destructive]:focus:ring-red-400 group-[.destructive]:focus:ring-offset-red-600", className),
  "toast-close": "",
  ...props
}, /* @__PURE__ */ React6.createElement(X, {
  className: "h-4 w-4"
})));
ToastClose.displayName = ToastPrimitives.Close.displayName;
var ToastTitle = /* @__PURE__ */ React6.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React6.createElement(ToastPrimitives.Title, {
  ref,
  className: cn("text-sm font-semibold", className),
  ...props
}));
ToastTitle.displayName = ToastPrimitives.Title.displayName;
var ToastDescription = /* @__PURE__ */ React6.forwardRef(({ className, ...props }, ref) => /* @__PURE__ */ React6.createElement(ToastPrimitives.Description, {
  ref,
  className: cn("text-sm opacity-90", className),
  ...props
}));
ToastDescription.displayName = ToastPrimitives.Description.displayName;

// src/components/ui/progress.tsx
import * as React10 from "react";

// ../../node_modules/.pnpm/@radix-ui+react-progress@1.1.7_@types+react-dom@19.0.0_@types+react@19.0.0_react-dom@19.0.0_react@19.0.0__react@19.0.0/node_modules/@radix-ui/react-progress/dist/index.mjs
import * as React9 from "react";

// ../../node_modules/.pnpm/@radix-ui+react-context@1.1.2_@types+react@19.0.0_react@19.0.0/node_modules/@radix-ui/react-context/dist/index.mjs
import * as React7 from "react";
import { jsx } from "react/jsx-runtime";
function createContextScope(scopeName, createContextScopeDeps = []) {
  let defaultContexts = [];
  function createContext3(rootComponentName, defaultContext) {
    const BaseContext = React7.createContext(defaultContext);
    const index = defaultContexts.length;
    defaultContexts = [...defaultContexts, defaultContext];
    const Provider2 = /* @__PURE__ */ __name((props) => {
      const { scope, children, ...context } = props;
      const Context = scope?.[scopeName]?.[index] || BaseContext;
      const value = React7.useMemo(() => context, Object.values(context));
      return /* @__PURE__ */ jsx(Context.Provider, { value, children });
    }, "Provider");
    Provider2.displayName = rootComponentName + "Provider";
    function useContext2(consumerName, scope) {
      const Context = scope?.[scopeName]?.[index] || BaseContext;
      const context = React7.useContext(Context);
      if (context) return context;
      if (defaultContext !== void 0) return defaultContext;
      throw new Error(`\`${consumerName}\` must be used within \`${rootComponentName}\``);
    }
    __name(useContext2, "useContext2");
    return [Provider2, useContext2];
  }
  __name(createContext3, "createContext3");
  const createScope = /* @__PURE__ */ __name(() => {
    const scopeContexts = defaultContexts.map((defaultContext) => {
      return React7.createContext(defaultContext);
    });
    return /* @__PURE__ */ __name(function useScope(scope) {
      const contexts = scope?.[scopeName] || scopeContexts;
      return React7.useMemo(
        () => ({ [`__scope${scopeName}`]: { ...scope, [scopeName]: contexts } }),
        [scope, contexts]
      );
    }, "useScope");
  }, "createScope");
  createScope.scopeName = scopeName;
  return [createContext3, composeContextScopes(createScope, ...createContextScopeDeps)];
}
__name(createContextScope, "createContextScope");
function composeContextScopes(...scopes) {
  const baseScope = scopes[0];
  if (scopes.length === 1) return baseScope;
  const createScope = /* @__PURE__ */ __name(() => {
    const scopeHooks = scopes.map((createScope2) => ({
      useScope: createScope2(),
      scopeName: createScope2.scopeName
    }));
    return /* @__PURE__ */ __name(function useComposedScopes(overrideScopes) {
      const nextScopes = scopeHooks.reduce((nextScopes2, { useScope, scopeName }) => {
        const scopeProps = useScope(overrideScopes);
        const currentScope = scopeProps[`__scope${scopeName}`];
        return { ...nextScopes2, ...currentScope };
      }, {});
      return React7.useMemo(() => ({ [`__scope${baseScope.scopeName}`]: nextScopes }), [nextScopes]);
    }, "useComposedScopes");
  }, "createScope");
  createScope.scopeName = baseScope.scopeName;
  return createScope;
}
__name(composeContextScopes, "composeContextScopes");

// ../../node_modules/.pnpm/@radix-ui+react-primitive@2.1.3_@types+react-dom@19.0.0_@types+react@19.0.0_react-dom@19.0.0_react@19.0.0__react@19.0.0/node_modules/@radix-ui/react-primitive/dist/index.mjs
import * as React8 from "react";
import * as ReactDOM from "react-dom";
import { createSlot } from "@radix-ui/react-slot";
import { jsx as jsx2 } from "react/jsx-runtime";
var NODES = [
  "a",
  "button",
  "div",
  "form",
  "h2",
  "h3",
  "img",
  "input",
  "label",
  "li",
  "nav",
  "ol",
  "p",
  "select",
  "span",
  "svg",
  "ul"
];
var Primitive = NODES.reduce((primitive, node) => {
  const Slot2 = createSlot(`Primitive.${node}`);
  const Node = React8.forwardRef((props, forwardedRef) => {
    const { asChild, ...primitiveProps } = props;
    const Comp = asChild ? Slot2 : node;
    if (typeof window !== "undefined") {
      window[Symbol.for("radix-ui")] = true;
    }
    return /* @__PURE__ */ jsx2(Comp, { ...primitiveProps, ref: forwardedRef });
  });
  Node.displayName = `Primitive.${node}`;
  return { ...primitive, [node]: Node };
}, {});

// ../../node_modules/.pnpm/@radix-ui+react-progress@1.1.7_@types+react-dom@19.0.0_@types+react@19.0.0_react-dom@19.0.0_react@19.0.0__react@19.0.0/node_modules/@radix-ui/react-progress/dist/index.mjs
import { jsx as jsx3 } from "react/jsx-runtime";
var PROGRESS_NAME = "Progress";
var DEFAULT_MAX = 100;
var [createProgressContext, createProgressScope] = createContextScope(PROGRESS_NAME);
var [ProgressProvider, useProgressContext] = createProgressContext(PROGRESS_NAME);
var Progress = React9.forwardRef(
  (props, forwardedRef) => {
    const {
      __scopeProgress,
      value: valueProp = null,
      max: maxProp,
      getValueLabel = defaultGetValueLabel,
      ...progressProps
    } = props;
    if ((maxProp || maxProp === 0) && !isValidMaxNumber(maxProp)) {
      console.error(getInvalidMaxError(`${maxProp}`, "Progress"));
    }
    const max = isValidMaxNumber(maxProp) ? maxProp : DEFAULT_MAX;
    if (valueProp !== null && !isValidValueNumber(valueProp, max)) {
      console.error(getInvalidValueError(`${valueProp}`, "Progress"));
    }
    const value = isValidValueNumber(valueProp, max) ? valueProp : null;
    const valueLabel = isNumber(value) ? getValueLabel(value, max) : void 0;
    return /* @__PURE__ */ jsx3(ProgressProvider, { scope: __scopeProgress, value, max, children: /* @__PURE__ */ jsx3(
      Primitive.div,
      {
        "aria-valuemax": max,
        "aria-valuemin": 0,
        "aria-valuenow": isNumber(value) ? value : void 0,
        "aria-valuetext": valueLabel,
        role: "progressbar",
        "data-state": getProgressState(value, max),
        "data-value": value ?? void 0,
        "data-max": max,
        ...progressProps,
        ref: forwardedRef
      }
    ) });
  }
);
Progress.displayName = PROGRESS_NAME;
var INDICATOR_NAME = "ProgressIndicator";
var ProgressIndicator = React9.forwardRef(
  (props, forwardedRef) => {
    const { __scopeProgress, ...indicatorProps } = props;
    const context = useProgressContext(INDICATOR_NAME, __scopeProgress);
    return /* @__PURE__ */ jsx3(
      Primitive.div,
      {
        "data-state": getProgressState(context.value, context.max),
        "data-value": context.value ?? void 0,
        "data-max": context.max,
        ...indicatorProps,
        ref: forwardedRef
      }
    );
  }
);
ProgressIndicator.displayName = INDICATOR_NAME;
function defaultGetValueLabel(value, max) {
  return `${Math.round(value / max * 100)}%`;
}
__name(defaultGetValueLabel, "defaultGetValueLabel");
function getProgressState(value, maxValue) {
  return value == null ? "indeterminate" : value === maxValue ? "complete" : "loading";
}
__name(getProgressState, "getProgressState");
function isNumber(value) {
  return typeof value === "number";
}
__name(isNumber, "isNumber");
function isValidMaxNumber(max) {
  return isNumber(max) && !isNaN(max) && max > 0;
}
__name(isValidMaxNumber, "isValidMaxNumber");
function isValidValueNumber(value, max) {
  return isNumber(value) && !isNaN(value) && value <= max && value >= 0;
}
__name(isValidValueNumber, "isValidValueNumber");
function getInvalidMaxError(propValue, componentName) {
  return `Invalid prop \`max\` of value \`${propValue}\` supplied to \`${componentName}\`. Only numbers greater than 0 are valid max values. Defaulting to \`${DEFAULT_MAX}\`.`;
}
__name(getInvalidMaxError, "getInvalidMaxError");
function getInvalidValueError(propValue, componentName) {
  return `Invalid prop \`value\` of value \`${propValue}\` supplied to \`${componentName}\`. The \`value\` prop must be:
  - a positive number
  - less than the value passed to \`max\` (or ${DEFAULT_MAX} if no \`max\` prop is set)
  - \`null\` or \`undefined\` if the progress is indeterminate.

Defaulting to \`null\`.`;
}
__name(getInvalidValueError, "getInvalidValueError");
var Root3 = Progress;
var Indicator = ProgressIndicator;

// src/components/ui/progress.tsx
var Progress2 = /* @__PURE__ */ React10.forwardRef(({ className, value, ...props }, ref) => /* @__PURE__ */ React10.createElement(Root3, {
  ref,
  className: cn("relative h-4 w-full overflow-hidden rounded-full bg-secondary", className),
  ...props
}, /* @__PURE__ */ React10.createElement(Indicator, {
  className: "h-full w-full flex-1 bg-primary transition-all",
  style: {
    transform: `translateX(-${100 - (value || 0)}%)`
  }
})));
Progress2.displayName = Root3.displayName;
export {
  Badge,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
  Input,
  Label,
  Progress2 as Progress,
  Toast,
  ToastAction,
  ToastClose,
  ToastDescription,
  ToastProvider,
  ToastTitle,
  ToastViewport,
  badgeVariants,
  buttonVariants,
  cn
};
