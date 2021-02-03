load("github.com/cirrus-templates/golang@v0.1.0", "lint_task", "test_task")

def main(ctx):
    return [lint_task(), test_task()]
