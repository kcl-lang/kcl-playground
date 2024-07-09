export const SHARE_QUERY_KEY = 's';

export const SNIPPETS = [
  {
    label: 'Kubernetes',
    value: `apiVersion = "apps/v1"
kind = "Deployment"
metadata = {
    name = "nginx"
    labels.app = "nginx"
}
spec = {
    replicas = 3
    selector.matchLabels = metadata.labels
    template.metadata.labels = metadata.labels
    template.spec.containers = [
        {
            name = metadata.name
            image = "nginx:1.14.2"
            ports = [{ containerPort = 80 }]
        }
    ]
}
`,
  },
  {
    label: 'Hello World!',
    value: `a = "Hello" + " " + "World!"
`,
  },

  {
    label: 'Computing integers',
    value: `x = (10 + 2) * 30 + 5
`,
  },

  {
    label: 'Debugging values',
    value: `print("Output values")
print(true, false, 100)
`,
  },

  {
    label: 'Conditionals',
    value: `if True:
    print("Hello")
else:
    print("unreachable")
`,
  },
  {
    label: 'List',
    value: `data = ["one", "two", "three"]`,
  },
  {
    label: 'Dict',
    value: `Config = {
    key = "value"    
}`,
  },
];
