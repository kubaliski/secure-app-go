import json, os, sys

def load_or_empty(path):
    try:
        with open(path) as f:
            content = f.read().strip()
            if not content:
                return {"results": []}
            return json.loads(content)
    except (json.JSONDecodeError, FileNotFoundError, IOError):
        return {"results": []}

def main():
    vuln_path = sys.argv[1] if len(sys.argv) > 1 else "reports/vuln-report.json"
    fixed_path = sys.argv[2] if len(sys.argv) > 2 else "reports/fixed-report.json"

    vuln = load_or_empty(vuln_path)
    fixed = load_or_empty(fixed_path)

    nvuln = len(vuln.get("results", []))
    nfixed = len(fixed.get("results", []))

    summary = "## SecureApp - Verificacion de Correcciones\n\n"
    summary += "### Version Vulnerable\n**Hallazgos: " + str(nvuln) + "**\n"

    if nvuln == 0:
        summary += "No se detectaron hallazgos (posible error en el escaneo)\n"
    else:
        for r in vuln.get("results", []):
            p = r.get("path", "")
            l = r.get("start", {}).get("line", "")
            m = r.get("extra", {}).get("message", "")[:80]
            summary += "- " + p + ":" + str(l) + " - " + m + "\n"

    summary += "\n### Version Corregida\n**Hallazgos: " + str(nfixed) + "**\n"
    if nfixed == 0:
        summary += "OK - Todas las vulnerabilidades han sido corregidas\n"
    else:
        summary += "ATENCION - Quedan hallazgos pendientes\n"
        for r in fixed.get("results", []):
            p = r.get("path", "")
            l = r.get("start", {}).get("line", "")
            m = r.get("extra", {}).get("message", "")[:80]
            summary += "- " + p + ":" + str(l) + " - " + m + "\n"

    step_summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if step_summary:
        with open(step_summary, "a") as out:
            out.write(summary)
    print(summary)

if __name__ == "__main__":
    main()
