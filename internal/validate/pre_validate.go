package validate

func Pre(dir string) Report {
	issues := make([]Issue, 0, 4)

	goModFile, issue := checkGoMod(dir)
	if issue != nil {
		issues = append(issues, *issue)
	}

	modulePath, issue := checkModule(goModFile)
	if issue != nil {
		issues = append(issues, *issue)
	}

	name, issue := checkName(modulePath)
	if issue != nil {
		issues = append(issues, *issue)
	}

	issue = checkWritable(dir)
	if issue != nil {
		issues = append(issues, *issue)
	}

	maxSev := maxSeverity(issues)

	return Report{
		Stage:      "pre",
		Verdict:    maxSev,
		Issues:     issues,
		ModulePath: modulePath,
		Name:       name,
	}
}
