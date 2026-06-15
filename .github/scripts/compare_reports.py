import json, os, sys

def main():
    vuln_path = sys.argv[1] if len(sys.argv) > 1 else "reports/vuln-report.json"
    fixed_path = sys.argv[2] if len(sys.argv) > 2 else "reports/fixed-report.json"

    with open(vuln_path) as f:
        vuln = json.load(f)
    with open(fixed_path) as f:
        fixed = json.load(f)

    nvuln = len(vuln.get("results", []))
    nfixed = len(fixed.get("results", []))

    summary = "## SecureApp - Verificacion de Correcciones\n\n"
    summary += "### Version Vulnerable\n**Hallazgos: " + str(nvuln) + "**\n"

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

    step_summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if step_summary:
        with open(step_summary, "a") as out:
            out.write(summary)
    print(summary)

if __name__ == "__main__":
    main()
