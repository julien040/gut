package controller

import (
	"fmt"
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Attributes(cmd *cobra.Command, args []string) error {
	selectedTemplate := selectAttributeTemplate()

	// Read the current attributes for the current directory
	exists := true
	currentAttributes, err := os.ReadFile(".gitattributes")
	switch {
	case os.IsNotExist(err):
		// File does not exist, create it
		exists = false
	case err != nil:
		// Some other error occurred
		exitOnError("I'm sorry, I can't read the .gitattributes file 😓. Please, can you check if I have the right permissions?", err)
	}

	alreadyIncluded := map[string]struct{}{}
	for _, line := range splitStringByNewLine(string(currentAttributes)) {
		if line == "" {
			continue
		}
		alreadyIncluded[line] = struct{}{}
	}

	toAdd := []string{}
	for _, line := range splitStringByNewLine(GitAttributes[selectedTemplate]) {
		if line == "" {
			continue
		}

		if _, ok := alreadyIncluded[line]; !ok {
			toAdd = append(toAdd, line)
		}
	}

	if len(toAdd) == 0 {
		fmt.Println("The selected template is already included in the .gitattributes file.")
		return nil
	}

	// Create the file
	if !exists {
		file, err := os.Create(".gitattributes")
		if err != nil {
			exitOnError("I'm sorry, I can't create the .gitattributes file 😓. Please, can you check if I have the right permissions?", err)
		}
		file.Close()
	}

	// Append the new attributes to the file
	file, err := os.OpenFile(".gitattributes", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		exitOnError("I'm sorry, I can't open the .gitattributes file 😓. Please, can you check if I have the right permissions?", err)
	}

	defer file.Close()
	file.WriteString("\n# " + selectedTemplate + " template downloaded with gut\n")
	for _, line := range toAdd {
		file.WriteString(line + "\n")
	}

	return nil
}

func selectAttributeTemplate() string {
	// Ask the user to select a gitignore template
	var templateNames []string
	for key := range GitAttributes {
		templateNames = append(templateNames, key)
	}
	// Sort the template names
	sort.Strings(templateNames)

	var selectedTemplate int
	prompt := &survey.Select{
		Message: "Select a gitignore template:",
		Options: templateNames,
	}
	err := survey.AskOne(prompt, &selectedTemplate)
	if err != nil {
		exitOnKnownError(errorReadInput, err)
	}
	return templateNames[selectedTemplate]
}

var GitAttributes = map[string]string{
	"R": `# Basic .gitattributes for a R repo.

# Source files
# ============
*.Rdata binary
*.RData binary
*.rda   binary
*.rdb   binary
*.rds   binary
*.Rd    text
*.Rdx   binary
*.Rmd	  text
*.R  	  text
*.Rproj text
*.[Rr]md   linguist-detectable`,
	"Common": `# Common settings that generally should always be used with your language specific settings

# Auto detect text files and perform LF normalization
*          text=auto

#
# The above will handle all files NOT found below
#

# Documents
*.bibtex   text diff=bibtex
*.doc      diff=astextplain
*.DOC      diff=astextplain
*.docx     diff=astextplain
*.DOCX     diff=astextplain
*.dot      diff=astextplain
*.DOT      diff=astextplain
*.pdf      diff=astextplain
*.PDF      diff=astextplain
*.rtf      diff=astextplain
*.RTF      diff=astextplain
*.md       text diff=markdown
*.mdx      text diff=markdown
*.tex      text diff=tex
*.adoc     text
*.textile  text
*.mustache text
*.csv      text eol=crlf
*.tab      text
*.tsv      text
*.txt      text
*.sql      text
*.epub     diff=astextplain

# Graphics
*.png      binary
*.jpg      binary
*.jpeg     binary
*.gif      binary
*.tif      binary
*.tiff     binary
*.ico      binary
# SVG treated as text by default.
*.svg      text
# If you want to treat it as binary,
# use the following line instead.
# *.svg    binary
*.eps      binary

# Scripts
*.bash     text eol=lf
*.fish     text eol=lf
*.ksh      text eol=lf
*.sh       text eol=lf
*.zsh      text eol=lf
# These are explicitly windows files and should use crlf
*.bat      text eol=crlf
*.cmd      text eol=crlf
*.ps1      text eol=crlf

# Serialisation
*.json     text
*.toml     text
*.xml      text
*.yaml     text
*.yml      text

# Archives
*.7z       binary
*.bz       binary
*.bz2      binary
*.bzip2    binary
*.gz       binary
*.lz       binary
*.lzma     binary
*.rar      binary
*.tar      binary
*.taz      binary
*.tbz      binary
*.tbz2     binary
*.tgz      binary
*.tlz      binary
*.txz      binary
*.xz       binary
*.Z        binary
*.zip      binary
*.zst      binary

# Text files where line endings should be preserved
*.patch    -text

#
# Exclude files from exporting
#

.gitattributes export-ignore
.gitignore     export-ignore
.gitkeep       export-ignore`,
	"PHP": `# Auto detect text files and perform LF normalization
*         text=auto

# PHP files
*.php     text eol=lf diff=php
*.phpt    text eol=lf diff=php
*.phtml   text eol=lf diff=html
*.twig    text eol=lf
*.phar    binary

# Configuration
phpcs.xml    text eol=lf
phpunit.xml  text eol=lf
phpstan.neon text eol=lf
psalm.xml    text eol=lf`,
	"Mathematica": `# Basic .gitattributes for a Mathematica repo.

# Source files
# Caution: *.m also matches Matlab files.
# ============
*.nb             text diff=mathematica
*.wls            text diff=mathematica
*.wl             text diff=mathematica
*.m              text diff=mathematica

# Test files
# ==========
*.wlt            text diff=mathematica
*.mt             text diff=mathematica

# Binary files
# ============
*.mx             binary`,
	"TinaCMS": `# https://tina.io/
.tina/__generated__/  linguist-generated`,
	"C++": `# Sources
*.c     text diff=cpp
*.cc    text diff=cpp
*.cxx   text diff=cpp
*.cpp   text diff=cpp
*.cpi   text diff=cpp
*.c++   text diff=cpp
*.hpp   text diff=cpp
*.h     text diff=cpp
*.h++   text diff=cpp
*.hh    text diff=cpp

# Compiled Object files
*.slo   binary
*.lo    binary
*.o     binary
*.obj   binary

# Precompiled Headers
*.gch   binary
*.pch   binary

# Compiled Dynamic libraries
*.so    binary
*.dylib binary
*.dll   binary

# Compiled Static libraries
*.lai   binary
*.la    binary
*.a     binary
*.lib   binary

# Executables
*.exe   binary
*.out   binary
*.app   binary`,
	"Python": `# Basic .gitattributes for a python repo.

# Source files
# ============
*.pxd    text diff=python
*.py     text diff=python
*.py3    text diff=python
*.pyw    text diff=python
*.pyx    text diff=python
*.pyz    text diff=python
*.pyi    text diff=python

# Binary files
# ============
*.db     binary
*.p      binary
*.pkl    binary
*.pickle binary
*.pyc    binary export-ignore
*.pyo    binary export-ignore
*.pyd    binary

# Jupyter notebook
*.ipynb  text eol=lf

# Note: .db, .p, and .pkl files are associated
# with the python modules ` + "`" + `` + "`" + `pickle` + "`" + `` + "`" + `, ` + "`" + `` + "`" + `dbm.*` + "`" + `` + "`" + `,
# ` + "`" + `` + "`" + `shelve` + "`" + `` + "`" + `, ` + "`" + `` + "`" + `marshal` + "`" + `` + "`" + `, ` + "`" + `` + "`" + `anydbm` + "`" + `` + "`" + `, & ` + "`" + `` + "`" + `bsddb` + "`" + `` + "`" + `
# (among others).`,
	"Fortran": `# Handle line endings automatically for files detected as text
# and leave all files detected as binary untouched.
*     text=auto

# Fortran files
*.f   text diff=fortran
*.for text diff=fortran
*.f90 text diff=fortran
*.f95 text diff=fortran
*.f03 text diff=fortran`,
	"PowerShell": `# Basic .gitattributes for a PowerShell repo.

# Source files
# ============
*.ps1      text eol=crlf
*.ps1x     text eol=crlf
*.psm1     text eol=crlf
*.psd1     text eol=crlf
*.ps1xml   text eol=crlf
*.pssc     text eol=crlf
*.psrc     text eol=crlf
*.cdxml    text eol=crlf`,
	"Rails": `# These settings are for Ruby On Rails project
#
#
# Define a dummy ours merge strategy with:
#
# $ git config --global merge.ours.driver true
schema.rb merge=ours diff=ruby`,
	"Delphi": `#
# Project files
#

# Delphi form module
*.dfm  text

# Delphi project options file
*.dof  text

# Desktop configuration
*.dsk  binary

#
# General
#

# Text file
*.txt  text

# Gettext Portable Object
*.po   binary

# Backup
*.bak  binary

# Config file
*.cfg  text

# Compiled Help File - contains html. See also package chm
*.chm  binary

# Comma Separated Values text file format
*.csv  text eol=crlf

# Directly executable program
*.exe  binary

# Help file
*.hlp  binary

# Initialization file
*.ini  text

# OpenDocument text document
*.odt  text

# Portable Document Format
*.pdf  text

# PostScript
*.ps   text

# Rich Text Format text file
*.rtf  text

#
# Image
#

# Portable network graphic
*.png  binary

# Bitmap
*.bmp  binary

# Icon
*.ico  binary

# Pixmap
*.xpm  binary

# Portable pixmap
*.ppm  binary

# Portable graymap
*.pgm  binary

# Portable bitmap
*.pbm  binary

# Lossy graphics file
*.jpg  binary

#
# XML
#

*.xml  text
*.dtd  text
*.xsd  text
*.xsl  text

#
# Web
#

# Hypertext Markup Language
*.html text diff=html

# Cascading style sheet
*.css  text diff=css

#
# Archive
#

# tape archive
*.tar  binary

# archive file
*.zip  binary

#
# Database
#

# Database file
*.dbf  binary

# Multiple index file
*.mdx  binary

#
# Spreadsheet
#

# OpenOffice.org Calc spreadsheet File Format
*.ods  binary

# Microsoft Excel File Format
*.xls  binary

# Microsoft Office Open XML Excel File Format
*.xlsx binary`,
	"Rust": `# Auto detect text files and perform normalization
*          text=auto

*.rs       text diff=rust
*.toml     text diff=toml
Cargo.lock text`,
	"Perl": `# Basic .gitattributes for a perl repo.

# Source files
# ============
*.pl text diff=perl
*.pm text diff=perl`,
	"Elixir": `# Handle line endings automatically for files detected as text
# and leave all files detected as binary untouched.
*     text=auto

# Elixir files
*.ex  diff=elixir
*.exs diff=elixir`,
	"DyalogAPL": `# Basic .gitattributes for a Dyalog APL repo.
# This template includes MiServer

# Source files
# ============
*.apl?          linguist-language=APL
*.dyalog        linguist-language=APL
*.dyapp         linguist-language=APL
*.mipage        linguist-language=APL

# Config and log files
# ====================
.dcfg           linguist-language=JSON5
.dlf            linguist-language=JSON5`,
	"Web": `## GITATTRIBUTES FOR WEB PROJECTS
#
# These settings are for any web project.
#
# Details per file setting:
#   text    These files should be normalized (i.e. convert CRLF to LF).
#   binary  These files are binary and should be left untouched.
#
# Note that binary is a macro for -text -diff.
######################################################################

# Auto detect
##   Handle line endings automatically for files detected as
##   text and leave all files detected as binary untouched.
##   This will handle all files NOT defined below.
*                 text=auto

# Source code
*.bash            text eol=lf
*.bat             text eol=crlf
*.cmd             text eol=crlf
*.coffee          text
*.css             text diff=css
*.htm             text diff=html
*.html            text diff=html
*.inc             text
*.ini             text
*.js              text
*.mjs             text
*.cjs             text
*.json            text
*.jsx             text
*.less            text
*.ls              text
*.map             text -diff
*.od              text
*.onlydata        text
*.php             text diff=php
*.pl              text
*.ps1             text eol=crlf
*.py              text diff=python
*.rb              text diff=ruby
*.sass            text
*.scm             text
*.scss            text diff=css
*.sh              text eol=lf
.husky/*          text eol=lf
*.sql             text
*.styl            text
*.tag             text
*.ts              text
*.tsx             text
*.xml             text
*.xhtml           text diff=html

# Docker
Dockerfile        text

# Documentation
*.ipynb           text eol=lf
*.markdown        text diff=markdown
*.md              text diff=markdown
*.mdwn            text diff=markdown
*.mdown           text diff=markdown
*.mkd             text diff=markdown
*.mkdn            text diff=markdown
*.mdtxt           text
*.mdtext          text
*.txt             text
AUTHORS           text
CHANGELOG         text
CHANGES           text
CONTRIBUTING      text
COPYING           text
copyright         text
*COPYRIGHT*       text
INSTALL           text
license           text
LICENSE           text
NEWS              text
readme            text
*README*          text
TODO              text

# Templates
*.dot             text
*.ejs             text
*.erb             text
*.haml            text
*.handlebars      text
*.hbs             text
*.hbt             text
*.jade            text
*.latte           text
*.mustache        text
*.njk             text
*.phtml           text
*.svelte          text
*.tmpl            text
*.tpl             text
*.twig            text
*.vue             text

# Configs
*.cnf             text
*.conf            text
*.config          text
.editorconfig     text
*.env             text
.gitattributes    text
.gitconfig        text
.htaccess         text
*.lock            text -diff
package.json      text eol=lf
package-lock.json text eol=lf -diff
pnpm-lock.yaml    text eol=lf -diff
.prettierrc       text
yarn.lock         text -diff
*.toml            text
*.yaml            text
*.yml             text
browserslist      text
Makefile          text
makefile          text
# Fixes syntax highlighting on GitHub to allow comments
tsconfig.json     linguist-language=JSON-with-Comments

# Heroku
Procfile          text

# Graphics
*.ai              binary
*.bmp             binary
*.eps             binary
*.gif             binary
*.gifv            binary
*.ico             binary
*.jng             binary
*.jp2             binary
*.jpg             binary
*.jpeg            binary
*.jpx             binary
*.jxr             binary
*.pdf             binary
*.png             binary
*.psb             binary
*.psd             binary
# SVG treated as an asset (binary) by default.
*.svg             text
# If you want to treat it as binary,
# use the following line instead.
# *.svg           binary
*.svgz            binary
*.tif             binary
*.tiff            binary
*.wbmp            binary
*.webp            binary

# Audio
*.kar             binary
*.m4a             binary
*.mid             binary
*.midi            binary
*.mp3             binary
*.ogg             binary
*.ra              binary

# Video
*.3gpp            binary
*.3gp             binary
*.as              binary
*.asf             binary
*.asx             binary
*.avi             binary
*.fla             binary
*.flv             binary
*.m4v             binary
*.mng             binary
*.mov             binary
*.mp4             binary
*.mpeg            binary
*.mpg             binary
*.ogv             binary
*.swc             binary
*.swf             binary
*.webm            binary

# Archives
*.7z              binary
*.gz              binary
*.jar             binary
*.rar             binary
*.tar             binary
*.zip             binary

# Fonts
*.ttf             binary
*.eot             binary
*.otf             binary
*.woff            binary
*.woff2           binary

# Executables
*.exe             binary
*.pyc             binary
# Prevents massive diffs caused by vendored, minified files
**/.yarn/releases/**   binary
**/.yarn/plugins/**    binary

# RC files (like .babelrc or .eslintrc)
*.*rc             text

# Ignore files (like .npmignore or .gitignore)
*.*ignore         text

# Prevents massive diffs from built files
dist/*            binary`,
	"ObjectiveC": `# compare .pbxproj files as binary and always merge as union
*.pbxproj binary -merge=union
*.m       text diff=objc`,
	"Lua": `# Basic .gitattributes for a Lua repo.

# Source files
# ============
*.lua       text

# Luadoc output
# =============
*.html      text diff=html
*.css       text diff=css`,
	"Markdown": `# Apply override to all files in the directory
*.md linguist-detectable`,
	"Java": `# Java sources
*.java          text diff=java
*.kt            text diff=kotlin
*.groovy        text diff=java
*.scala         text diff=java
*.gradle        text diff=java
*.gradle.kts    text diff=kotlin

# These files are text and should be normalized (Convert crlf => lf)
*.css           text diff=css
*.scss          text diff=css
*.sass          text
*.df            text
*.htm           text diff=html
*.html          text diff=html
*.js            text
*.mjs           text
*.cjs           text
*.jsp           text
*.jspf          text
*.jspx          text
*.properties    text
*.tld           text
*.tag           text
*.tagx          text
*.xml           text

# These files are binary and should be left untouched
# (binary is a macro for -text -diff)
*.class         binary
*.dll           binary
*.ear           binary
*.jar           binary
*.so            binary
*.war           binary
*.jks           binary

# Common build-tool wrapper scripts ('.cmd' versions are handled by 'Common.gitattributes')
mvnw            text eol=lf
gradlew         text eol=lf

# These are explicitly windows files and should use crlf
*.bat           text eol=crlf`,
	"Servoy": `# Auto detect text files and perform LF normalization
*          text=auto

*.frm -text
*.val -text
*.tbl -text
*.rel -text
*.obj -text
*.dbi -text
*.sec -text
*.css text diff=css
*.js       eol=lf
*.mjs      eol=lf
*.cjs      eol=lf`,
	"CSharp": `# Auto detect text files and perform LF normalization
*          text=auto

*.cs       text diff=csharp
*.cshtml   text diff=html
*.csx      text diff=csharp
*.sln      text eol=crlf
*.csproj   text eol=crlf`,
	"Vim": `# Basic .gitattributes for a Vim repo.
# Vim on Linux works with LF only, Vim on Windows works with both LF and CRLF

# Source files
# ============
*.vim text eol=lf
.vimrc text eol=lf
.gvimrc text eol=lf`,
	"Ada": `# Handle line endings automatically for files detected as text
# and leave all files detected as binary untouched.
*     text=auto

# Ada files
*.ada text diff=ada
*.adb text diff=ada
*.ads text diff=ada

# GNAT Project files (syntax similar to Ada)
*.gpr linguist-language=Ada`,
	"Pascal": `# Lazarus Project Information file (stored in XML; contains project-specific settings)
*.lpi      text

# Lazarus Program file; contains Pascal source of main program
*.lpr      text

# Lazarus Form file; contains configuration information for all objects on a form
# (stored in a Lazarus-specific format; the actions are described by Pascal source code in a corresponding *.pas file)
*.lfm      text

# Unit with Pascal code (typically for a form stored in a corresponding *.lfm file)
*.pas      text diff=pascal

# Pascal code
*.p        text diff=pascal
*.pp       text diff=pascal

# Lazarus Resource file (this is a generated file; not to be confused with a Windows resource file).
*.lrs      text

# Compiled unit, symbols part.
*.ppu      binary

# (1) Compiled unit, code part.
# (2) Compiled code from other compilers (e.g. gcc)
*.o        binary

# Object resource, automatically generated from {$R} directive.
*.or       text

# (1) Compiled unit, code part for smartlinking (on some platforms)
# (2) Compiled code from other compilers (e.g. gcc) linked into a static library
*.a        binary

# Lazarus package information file. (stored in XML; contains package-specific settings)
*.lpk      text

# Include file
*.inc      text

# Lazarus Project Session file. See Project Options -> Save session information in
*.lps      text

# Lazarus Resourcestring table created when saving a lfm file and i18n is enabled. It contains the TTranslateString properties of the lfm.
*.lrt      text

# Resourcestring table created by the compiler for every unit with a resourcestring section.
*.rst      text

# Resourcestring table in JSON format created by FPC 2.7.1 for units with resourcestring section.
*.rsj      text

# Compilation session for a project
*.compiled text

# Resource file
*.res      text

# Lazarus resource form file
*.rc       text

# Icon file
*.ico      binary`,
	"MicrosoftShell": `
# Basic .gitattributes for a Microsoft Shell repo.

# Source files
# ============
*.msh      text eol=crlf
*.msh1     text eol=crlf
*.msh2     text eol=crlf
*.mshxml   text eol=crlf
*.msh1xml  text eol=crlf
*.msh2xml  text eol=crlf
*.mcf      text eol=crlf`,
	"Drupal": `# Drupal git normalization
# @see https://www.kernel.org/pub/software/scm/git/docs/gitattributes.html
# @see https://www.drupal.org/node/1542048

# Normally these settings would be done with macro attributes for improved
# readability and easier maintenance. However macros can only be defined at the
# repository root directory. Drupal avoids making any assumptions about where it
# is installed.

# Define text file attributes.
# - Treat them as text.
# - Ensure no CRLF line-endings, neither on checkout nor on checkin.
# - Detect whitespace errors.
#   - Exposed by default in ` + "`" + `git diff --color` + "`" + ` on the CLI.
#   - Validate with ` + "`" + `git diff --check` + "`" + `.
#   - Deny applying with ` + "`" + `git apply --whitespace=error-all` + "`" + `.
#   - Fix automatically with ` + "`" + `git apply --whitespace=fix` + "`" + `.

*.bash    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.config  text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.css     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.dist    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.engine  text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.html    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=html
*.inc     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.install text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.js      text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.mjs     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.cjs     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.json    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.lock    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.map     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.md      text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.module  text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.php     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.po      text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.profile text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.script  text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.sh      text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.sql     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.svg     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.theme   text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2 diff=php
*.twig    text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.txt     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.xml     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2
*.yml     text eol=lf whitespace=blank-at-eol,-blank-at-eof,-space-before-tab,tab-in-indent,tabwidth=2

# Define binary file attributes.
# - Do not treat them as text.
# - Include binary diff in patches instead of "binary files differ."
*.eot     -text diff
*.exe     -text diff
*.gif     -text diff
*.gz      -text diff
*.ico     -text diff
*.jpeg    -text diff
*.jpg     -text diff
*.otf     -text diff
*.phar    -text diff
*.png     -text diff
*.svgz    -text diff
*.ttf     -text diff
*.woff    -text diff
*.woff2   -text diff`,
	"Flutter": `# Auto detect text files and perform LF normalization
*                 text=auto

# Always perform LF normalization
*.dart            text
*.gradle          text
*.html            text
*.java            text
*.json            text
*.md              text
*.py              text
*.sh              text
*.txt             text
*.xml             text
*.yaml            text

# Make sure that these Windows files always have CRLF line endings at checkout
*.bat             text eol=crlf
*.ps1             text eol=crlf
*.rc              text eol=crlf
*.sln             text eol=crlf
*.props           text eol=crlf
*.vcxproj         text eol=crlf
*.vcxproj.filters text eol=crlf
# Including templates
*.sln.tmpl        text eol=crlf
*.props.tmpl      text eol=crlf
*.vcxproj.tmpl    text eol=crlf

# Never perform LF normalization
*.ico             binary
*.jar             binary
*.png             binary
*.zip             binary
*.ttf             binary
*.otf             binary`,
	"Hashicorp": `# HachiCorp Configuration Language
*.hcl eol=lf

# Terraform
*.tf      text eol=lf
*.tf.json text eol=lf
*.tfvars  text eol=lf

# Vagrant
Vagrantfile text eol=lf`,
	"sql": `# Basic .gitattributes for sql files

*.sql linguist-detectable=true
*.sql linguist-language=sql`,
	"FSharp": `# Auto detect text files and perform LF normalization
*          text=auto

*.fs       text diff=fsharp
*.fsx      text diff=fsharp
*.sln      text eol=crlf
*.fsproj   text eol=crlf`,
	"Ballerina": `# Handle line endings automatically for files detected as text
# and leave all files detected as binary untouched.
*               text=auto

#
# The above will handle all files NOT found below
#
# These files are text and should be normalized (Convert crlf => lf)
*.bal           text
*.bash          text eol=lf
*.css           text diff=css
*.df            text
*.htm           text diff=html
*.html          text diff=html
*.js            text
*.mjs           text
*.cjs           text
*.json          text
*.properties    text
*.sh            text eol=lf
*.tld           text
*.txt           text
*.tag           text
*.tagx          text
*.xml           text
*.yml           text

# These files are binary and should be left untouched
# (binary is a macro for -text -diff)
*.balx          binary
*.dll           binary
*.ear           binary
*.gif           binary
*.ico           binary
*.jpg           binary
*.jpeg          binary
*.png           binary
*.so            binary`,
	"Fountain": `# Handle line endings automatically for files detected as text
# and leave all files detected as binary untouched.
*          text=auto

# Fountain files
*.fountain text diff=fountain
*.pdf      binary
*.fdx      text`,
	"Unity": `# Define macros (only works in top-level gitattributes files)
[attr]lfs               filter=lfs diff=lfs merge=lfs -text
[attr]unity-json        eol=lf linguist-language=json
[attr]unity-yaml        merge=unityyamlmerge eol=lf linguist-language=yaml

# Optionally collapse Unity-generated files on GitHub diffs
# [attr]unity-yaml        merge=unityyamlmerge text linguist-language=yaml linguist-generated

# Unity source files
*.cginc                 text
*.compute               text linguist-language=hlsl
*.cs                    text diff=csharp
*.hlsl                  text linguist-language=hlsl
*.raytrace              text linguist-language=hlsl
*.shader                text

# Unity JSON files
*.asmdef                unity-json
*.asmref                unity-json
*.index                 unity-json
*.inputactions          unity-json
*.shadergraph           unity-json
*.shadersubgraph        unity-json

# Unity UI Toolkit files
*.tss                   text diff=css linguist-language=css
*.uss                   text diff=css linguist-language=css
*.uxml                  text linguist-language=xml linguist-detectable

# Unity YAML
*.anim                  unity-yaml
*.asset                 unity-yaml
*.brush                 unity-yaml
*.controller            unity-yaml
*.flare                 unity-yaml
*.fontsettings          unity-yaml
*.giparams              unity-yaml
*.guiskin               unity-yaml
*.lighting              unity-yaml
*.mask                  unity-yaml
*.mat                   unity-yaml
*.meta                  unity-yaml
*.mixer                 unity-yaml
*.overrideController    unity-yaml
*.playable              unity-yaml
*.prefab                unity-yaml
*.preset                unity-yaml
*.renderTexture         unity-yaml
*.scenetemplate         unity-yaml
*.shadervariants        unity-yaml
*.signal                unity-yaml
*.spriteatlas           unity-yaml
*.spriteatlasv2         unity-yaml
*.terrainlayer          unity-yaml
*.unity                 unity-yaml

# "physic" for 3D but "physics" for 2D
*.physicMaterial        unity-yaml
*.physicsMaterial2D     unity-yaml

# Exclude third-party plugins from GitHub stats
Assets/Plugins/**       linguist-vendored

# Unity LFS
*.cubemap               lfs
*.unitypackage          lfs

# 3D models
*.3dm                   lfs
*.3ds                   lfs
*.blend                 lfs
*.c4d                   lfs
*.collada               lfs
*.dae                   lfs
*.dxf                   lfs
*.FBX                   lfs
*.fbx                   lfs
*.jas                   lfs
*.lws                   lfs
*.lxo                   lfs
*.ma                    lfs
*.max                   lfs
*.mb                    lfs
*.obj                   lfs
*.ply                   lfs
*.skp                   lfs
*.stl                   lfs
*.ztl                   lfs

# Audio
*.aif                   lfs
*.aiff                  lfs
*.it                    lfs
*.mod                   lfs
*.mp3                   lfs
*.ogg                   lfs
*.s3m                   lfs
*.wav                   lfs
*.xm                    lfs

# Video
*.asf                   lfs
*.avi                   lfs
*.flv                   lfs
*.mov                   lfs
*.mp4                   lfs
*.mpeg                  lfs
*.mpg                   lfs
*.ogv                   lfs
*.wmv                   lfs

# Images
*.bmp                   lfs
*.exr                   lfs
*.gif                   lfs
*.hdr                   lfs
*.iff                   lfs
*.jpeg                  lfs
*.jpg                   lfs
*.pict                  lfs
*.png                   lfs
*.psd                   lfs
*.tga                   lfs
*.tif                   lfs
*.tiff                  lfs
*.webp                  lfs

# Compressed Archive
*.7z                    lfs
*.bz2                   lfs
*.gz                    lfs
*.rar                   lfs
*.tar                   lfs
*.zip                   lfs

# Compiled Dynamic Library
*.dll                   lfs
*.pdb                   lfs
*.so                    lfs

# Fonts
*.otf                   lfs
*.ttf                   lfs

# Executable/Installer
*.apk                   lfs
*.exe                   lfs

# Documents
*.pdf                   lfs

# ETC
*.a                     lfs
*.reason                lfs
*.rns                   lfs

# Spine export file for Unity
*.skel.bytes            lfs

# Exceptions for .asset files such as lightning pre-baking
LightingData.asset     binary`,
	"ActionScript": `# Adobe Flash authoring file
*.fla  text

# ActionScript file
*.as   text

# Flash XML file
*.xml  text

# Treat .swf and .swc as binary
# https://stackoverflow.com/q/1529178

# Compiled Flash file
*.swf  binary -crlf -diff -merge

# Compiled Flash Library/Script
*.swc  binary -crlf -diff -merge

# ActionScript Communication file
*.asc  text

# Flash JavaScript file
*.jsfl text`,
	"VisualStudio": `###############################################################################
# Set default behavior to automatically normalize line endings.
###############################################################################
*            text=auto

###############################################################################
# Set the merge driver for project and solution files
#
# Merging from the command prompt will add diff markers to the files if there
# are conflicts (Merging from VS is not affected by the settings below, in VS
# the diff markers are never inserted). Diff markers may cause the following
# file extensions to fail to load in VS. An alternative would be to treat
# these files as binary and thus will always conflict and require user
# intervention with every merge. To do so, just comment the entries below and
# uncomment the group further below
###############################################################################

*.sln        text eol=crlf
*.csproj     text eol=crlf
*.vbproj     text eol=crlf
*.vcxproj    text eol=crlf
*.vcproj     text eol=crlf
*.dbproj     text eol=crlf
*.fsproj     text eol=crlf
*.lsproj     text eol=crlf
*.wixproj    text eol=crlf
*.modelproj  text eol=crlf
*.sqlproj    text eol=crlf
*.wwaproj    text eol=crlf

*.xproj      text eol=crlf
*.props      text eol=crlf
*.filters    text eol=crlf
*.vcxitems   text eol=crlf


#*.sln       merge=binary
#*.csproj    merge=binary
#*.vbproj    merge=binary
#*.vcxproj   merge=binary
#*.vcproj    merge=binary
#*.dbproj    merge=binary
#*.fsproj    merge=binary
#*.lsproj    merge=binary
#*.wixproj   merge=binary
#*.modelproj merge=binary
#*.sqlproj   merge=binary
#*.wwaproj   merge=binary

#*.xproj     merge=binary
#*.props     merge=binary
#*.filters   merge=binary
#*.vcxitems  merge=binary`,
	"DevContainer": `# Fix syntax highlighting on GitHub to allow comments
.devcontainer.json linguist-language=JSON-with-Comments
devcontainer.json linguist-language=JSON-with-Comments`,
	"VisualStudioCode": `# Fix syntax highlighting on GitHub to allow comments
.vscode/*.json linguist-language=JSON-with-Comments`,
	"Go": `# Treat all Go files in this repo as binary, with no git magic updating
# line endings. Windows users contributing to Go will need to use a
# modern version of git and editors capable of LF line endings.

*.go -text diff=golang`,
}
