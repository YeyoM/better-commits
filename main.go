package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "time"

  tea "github.com/charmbracelet/bubbletea"

  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/textinput"
  "github.com/charmbracelet/bubbles/textarea"
  "github.com/charmbracelet/bubbles/timer"

  "github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
const timeout = 5 * time.Second


// the structure of the conventional commits are as follows:
// <type>[optional scope]: <description>
// [optional body]
// [optional footer]

type item struct {
  title, desc string
}

func (i item) Title() string {
  return i.title
}

func (i item) Description() string {
  return i.desc 
}

func (i item) FilterValue() string {
  return i.title 
}

type model struct {
  canCommit bool
  noChanges bool
  changesNotStaged bool
  finished bool
  timer timer.Model
  currentStep int
  commitType string
  typeOptions list.Model
  scope string
  scopeInput textinput.Model
  gitmoji string
  gitmojiOptions list.Model
  description string
  descriptionInput textinput.Model
  body string
  bodyInput textarea.Model
  footer string
  footerInput textarea.Model
}

func initialModel() model {
  typeOptions := []list.Item{
    item{title: "feat", desc: "A new feature"},
    item{title: "fix", desc: "A bug fix"},
    item{title: "docs", desc: "Documentation only changes"},
    item{title: "style", desc: "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)"},
    item{title: "refactor", desc: "A code change that neither fixes a bug nor adds a feature"},
    item{title: "perf", desc: "A code change that improves performance"},
    item{title: "test", desc: "Adding missing tests or correcting existing tests"},
    item{title: "build", desc: "Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
    item{title: "ci", desc: "Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)"},
    item{title: "chore", desc: "Other changes that don't modify src or test files"},
  }

  gitmojiOptions := []list.Item{
    item{title: "none", desc: "No emoji"},
    item{title: ":art:", desc: "Improving structure / format of the code"},
    item{title: ":zap:", desc: "Improving performance"},
    item{title: ":fire:", desc: "Removing code or files"},
    item{title: ":bug:", desc: "Fixing a bug"},
    item{title: ":ambulance:", desc: "Critical hotfix"},
    item{title: ":sparkles:", desc: "Introducing new features"},
    item{title: ":pencil:", desc: "Writing docs"},
    item{title: ":rocket:", desc: "Deploying stuff"},
    item{title: ":lipstick:", desc: "Updating the UI and style files"},
    item{title: ":tada:", desc: "Initial commit"},
    item{title: ":white_check_mark:", desc: "Adding tests"},
    item{title: ":lock:", desc: "Fixing security issues"},
    item{title: ":apple:", desc: "Fixing something on macOS"},
    item{title: ":penguin:", desc: "Fixing something on Linux"},
    item{title: ":checkered_flag:", desc: "Fixing something on Windows"},
    item{title: ":robot:", desc: "Fixing something on Android"},
    item{title: ":green_apple:", desc: "Fixing something on iOS"},
    item{title: ":bookmark:", desc: "Releasing / Version tags"},
    item{title: ":rotating_light:", desc: "Removing linter warnings"},
    item{title: ":construction:", desc: "Work in progress"},
    item{title: ":green_heart:", desc: "Fixing CI Build"},
    item{title: ":arrow_down:", desc: "Downgrading dependencies"},
    item{title: ":arrow_up:", desc: "Upgrading dependencies"},
    item{title: ":pushpin:", desc: "Pinning dependencies to specific versions"},
    item{title: ":construction_worker:", desc: "Adding CI build system"},
    item{title: ":chart_with_upwards_trend:", desc: "Adding analytics or tracking code"},
    item{title: ":recycle:", desc: "Refactoring code"},
    item{title: ":whale:", desc: "Work about Docker"},
    item{title: ":heavy_plus_sign:", desc: "Adding a dependency"},
    item{title: ":heavy_minus_sign:", desc: "Removing a dependency"},
    item{title: ":wrench:", desc: "Changing configuration files"},
    item{title: ":globe_with_meridians:", desc: "Internationalization and localization"},
    item{title: ":pencil2:", desc: "Fixing typos"},
    item{title: ":poop:", desc: "Writing bad code that needs to be improved"},
    item{title: ":rewind:", desc: "Reverting changes"},
    item{title: ":twisted_rightwards_arrows:", desc: "Merging branches"},
    item{title: ":package:", desc: "Updating compiled files or packages"},
    item{title: ":alien:", desc: "Updating code due to external API changes"},
    item{title: ":truck:", desc: "Moving or renaming files"},
    item{title: ":page_facing_up:", desc: "Adding or updating license"},
    item{title: ":boom:", desc: "Introducing breaking changes"},
    item{title: ":bento:", desc: "Adding or updating assets"},
    item{title: ":ok_hand:", desc: "Updating code due to code review changes"},
    item{title: ":wheelchair:", desc: "Improving accessibility"},
    item{title: ":bulb:", desc: "Documenting source code"},
    item{title: ":beers:", desc: "Writing code drunkenly"},
    item{title: ":speech_balloon:", desc: "Updating title and literals"},
    item{title: ":card_file_box:", desc: "Performing database related changes"},
    item{title: ":loud_sound:", desc: "Adding logs"},
    item{title: ":mute:", desc: "Removing logs"},
    item{title: ":busts_in_silhouette:", desc: "Adding contributor(s)"},
    item{title: ":children_crossing:", desc: "Improving user experience / usability"},
    item{title: ":building_construction:", desc: "Making architectural changes"},
    item{title: ":iphone:", desc: "Working on responsive design"},
    item{title: ":clown_face:", desc: "Mocking things"},
    item{title: ":egg:", desc: "Adding an easter egg"},
    item{title: ":see_no_evil:", desc: "Adding or updating a .gitignore file"},
    item{title: ":camera_flash:", desc: "Adding or updating snapshots"},
    item{title: ":alembic:", desc: "Experimenting new things"},
    item{title: ":mag:", desc: "Improving SEO"},
    item{title: ":wheel_of_dharma:", desc: "Work about Kubernetes"},
    item{title: ":label:", desc: "Adding or updating types (Flow, TypeScript)"},
    item{title: ":seedling:", desc: "Adding or updating seed files"},
    item{title: ":triangular_flag_on_post:", desc: "Adding, updating, or removing feature flags"},
    item{title: ":goal_net:", desc: "Catching errors"},
    item{title: ":dizzy:", desc: "Adding or updating animations and transitions"},
    item{title: ":wastebasket:", desc: "Deprecating code that needs to be cleaned up"},
    item{title: ":passport_control:", desc: "Working on sign up and sign in"},
  }

  model := model{
    canCommit: false,
    noChanges: false,
    changesNotStaged: false,
    finished: false,
    timer: timer.NewWithInterval(timeout, time.Second),
    currentStep: -1,
    commitType: "",
    typeOptions: list.New(typeOptions, list.NewDefaultDelegate(), 0, 0),
    scope: "",
    scopeInput: textinput.New(),
    gitmoji: "",
    gitmojiOptions: list.New(gitmojiOptions, list.NewDefaultDelegate(), 0, 0),
    description: "",
    descriptionInput: textinput.New(),
    body: "",
    bodyInput: textarea.New(),
    footer: "",
    footerInput: textarea.New(),
  }


  model.typeOptions.Title = "Select the type of change"

  model.scopeInput.Placeholder = "Scope (optional), press enter to continue"
  model.scopeInput.Focus()
  model.scopeInput.CharLimit = 50
  model.scopeInput.Width = 50

  model.gitmojiOptions.Title = "Select a gitmoji for the commit"

  model.descriptionInput.Placeholder = "Write a short, imperative tense description of the change"
  model.descriptionInput.Focus()
  model.descriptionInput.CharLimit = 50
  model.descriptionInput.Width = 50


  model.bodyInput.Placeholder = "Provide a longer description of the change (optional)"
  model.bodyInput.SetWidth(50)
  model.bodyInput.SetHeight(6)
  model.bodyInput.Focus()
  model.bodyInput.ShowLineNumbers = false
  model.bodyInput.FocusedStyle.CursorLine = lipgloss.NewStyle()

  model.footerInput.Placeholder = "List any breaking changes or issues closed by this change (optional)"
  model.footerInput.SetWidth(50)
  model.footerInput.SetHeight(6)
  model.footerInput.Focus()
  model.footerInput.ShowLineNumbers = false
  model.footerInput.FocusedStyle.CursorLine = lipgloss.NewStyle()

  return model
}

func (m model) Init() tea.Cmd {
  return m.timer.Init()
} 

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
    case timer.TickMsg:
      m.checkNoChanges()
      m.timer, cmd = m.timer.Update(msg)
      return m, cmd
    
    case timer.TimeoutMsg:
      if m.canCommit == false && m.noChanges == true {
        return m, tea.Quit
      }

    case tea.KeyMsg:
      if msg.String() == "ctrl+c" {
        m.commitType = ""
        m.scope = ""
        m.gitmoji = ""
        m.description = ""
        m.body = ""
        m.footer = ""
        return m, tea.Quit
      }

      if msg.String() == "enter" {
        switch m.currentStep {
          case -1: // check if there are changes to commit
            if m.canCommit && !m.noChanges {
              m.currentStep++
            }
            return m, nil
          case 0: // step 1 -> select the type of change (required)
            m.currentStep++
            if m.commitType == "" {
              m.commitType = "feat"
            }
            return m, nil
          case 1: // step 2 -> what is the scope of this change (optional)
            m.currentStep++
            return m, nil
          case 2: // step 3 -> select a gitmoji for the commit (optional)
            m.currentStep++
            return m, nil
          case 3: // step 4 -> write a short, imperative tense description of the change (required)
            if m.description == "" {
              docStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Margin(1, 2) // change the color to red
              return m, nil
            }
            docStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff")).Margin(1, 2) // change the color back to white
            m.currentStep++
            return m, nil
          case 4: // step 5 -> provide a longer description of the change (optional)
            m.currentStep++
            return m, nil
          case 5: // step 6 -> list any breaking changes or issues closed by this change (optional)
            m.currentStep++
            return m, nil
          case 6: // step 7 -> display the commit message (preview)
            m.CommitChanges() 
            time.Sleep(2 * time.Second)
            m.currentStep++
            return m, nil
          default:
            return m, tea.Quit
        }
      }

      // move to the previous step 
      if msg.String() == "tab" {
        switch m.currentStep {
          case 1:
            m.currentStep--
            return m, nil
          case 2:
            m.currentStep--
            return m, nil
          case 3:
            m.currentStep--
            return m, nil
          case 4:
            m.currentStep--
            return m, nil
          case 5:
            m.currentStep--
            return m, nil
          case 6:
            m.currentStep--
            return m, nil
        }
      }

    case tea.WindowSizeMsg:
      h, v := docStyle.GetFrameSize()
      m.typeOptions.SetSize(msg.Width - h, msg.Height - v - 4)
      m.scopeInput.Width = msg.Width - h
      m.gitmojiOptions.SetSize(msg.Width - h, msg.Height - v - 4)
      m.descriptionInput.Width = msg.Width - h
  }

  if m.currentStep == -1 {
    m.timer, cmd = m.timer.Update(msg)
    if !m.noChanges {
      m.timer.Toggle()
      return m, cmd
    }
  }

  // update the model when selecting an item from the list (typeOptions)
  if m.currentStep == 0 {
    m.typeOptions, cmd = m.typeOptions.Update(msg)
    m.commitType = string(m.typeOptions.SelectedItem().FilterValue())
  }

  // update the model when typing to the text input (scope)
  if m.currentStep == 1 {
    m.scopeInput, cmd = m.scopeInput.Update(msg)
    m.scope = m.scopeInput.Value()
  }

  // update the model when selecting an item from the list (gitmojiOptions)
  if m.currentStep == 2 {
    m.gitmojiOptions, cmd = m.gitmojiOptions.Update(msg)
    var possible = string(m.gitmojiOptions.SelectedItem().FilterValue())
    if possible == "none" {
      m.gitmoji = ""
    } else {
      m.gitmoji = possible
    }
  }

  // update the model when typing to the text input (description)
  if m.currentStep == 3 {
    m.descriptionInput, cmd = m.descriptionInput.Update(msg)
    m.description = m.descriptionInput.Value()
  }

  // update the model when typing to the text input (body)
  if m.currentStep == 4 {
    m.bodyInput, cmd = m.bodyInput.Update(msg)
    m.body = m.bodyInput.Value()
  }

  // update the model when typing to the text input (footer)
  if m.currentStep == 5 {
    m.footerInput, cmd = m.footerInput.Update(msg)
    m.footer = m.footerInput.Value()
  }

  return m, cmd
}

func (m model) HelpMenu() string {
    help := `
    Press Ctrl+C to quit at any time.
    Press Enter to proceed to the next step.
    Press Tab to go back to the previous step.
    `
    return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(help)
}

func (m model) EndHelpMenu() string {
    help := `
    Press Ctrl+C to quit 
    Press Enter to continue
    `
    return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(help)
}

func displayGitStatus() string {
  cmd := exec.Command("git", "status")
  cmd.Dir = "." // Set the command's working directory to the current directory

  output, err := cmd.Output()
  if err!= nil {
    return fmt.Sprintf("Error executing 'git status': %v", err)
  }

  status := string(output)
  formattedStatus := strings.TrimSpace(status)

  // Remove the fixed height and let the content expand
  formattedStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("#00A0FF")).Render(formattedStatus)

  formattedStatus = fmt.Sprintf("Current Git status:\n\n%s", formattedStatus)
  return formattedStatus
}

func (m *model) checkNoChanges() {
  gitStatus := displayGitStatus()
  if strings.Contains(gitStatus, "no changes added to commit") ||
    strings.Contains(gitStatus, "nothing to commit, working tree clean") {
    m.noChanges = true
    m.canCommit = false
  } else {
    m.noChanges = false
    m.canCommit = true
  }
}

func (m *model) checkChangesNotStaged() {
  gitStatus := displayGitStatus()
  if strings.Contains(gitStatus, "Changes not staged for commit") {
    m.changesNotStaged = true
  } 
}

func (m model) View() string {
  gitStatus := displayGitStatus()

  switch m.currentStep {
    case -1: 
      m.checkNoChanges()
      m.checkChangesNotStaged()
      if m.noChanges {
        s := m.timer.View()
        return docStyle.Render(gitStatus + "\n\n" + "No changes to commit. Exiting in " + s + " seconds... " + "\n")
      } else if m.changesNotStaged {
        return docStyle.Render(gitStatus + "\n\n" + "There are some changes that are not staged to commit, be sure before proceeding." + "\n" + m.HelpMenu())
      }
      return docStyle.Render(gitStatus + "\n" + m.HelpMenu())
    case 0: // select the type of change
      return docStyle.Render(m.typeOptions.View() + "\n" + m.HelpMenu())
    case 1: // what is the scope of this change
      return docStyle.Render(m.scopeInput.View() + "\n" + m.HelpMenu())
    case 2: // select a gitmoji for the commit 
      return docStyle.Render(m.gitmojiOptions.View() + "\n" + m.HelpMenu())
    case 3: // write a short, imperative tense description of the change 
      return docStyle.Render(m.descriptionInput.View() + "\n" + m.HelpMenu())
    case 4: // provide a longer description of the change 
      return docStyle.Render(m.bodyInput.View() + "\n" + m.HelpMenu())
    case 5: // list any breaking changes or issues closed by this change 
      return docStyle.Render(m.footerInput.View() + "\n" + m.HelpMenu())
    case 6: // display the commit message (preview)
      return docStyle.Render(m.PrintCommitMessage() + "\n" + m.EndHelpMenu())
    default:
      return "Exiting..."
  }

  return ""
}
func (m model) PlainCommitMessage() string {
  if m.scope == "" {
    return fmt.Sprintf("%s: %s %s\n\n%s\n\n%s",
      m.commitType, m.gitmoji, m.description, m.body, m.footer)
  }

  return fmt.Sprintf("%s[%s]: %s %s\n\n%s\n\n%s",
    m.commitType, m.scope, m.gitmoji, m.description, m.body, m.footer)
}

func (m model) PrintCommitMessage() string {
  commitTypeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
  scopeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
  gitmojiStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
  descriptionStyle := lipgloss.NewStyle().Bold(true)
  bodyStyle := lipgloss.NewStyle().MarginTop(1)
  footerStyle := lipgloss.NewStyle().MarginTop(1).Italic(true)
  confirmStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)

  formattedCommitType := commitTypeStyle.Render(m.commitType)
  formattedScope := scopeStyle.Render(m.scope)
  formattedGitmoji := gitmojiStyle.Render(m.gitmoji)
  formattedDescription := descriptionStyle.Render(m.description)
  formattedBody := bodyStyle.Render(m.body)
  formattedFooter := footerStyle.Render(m.footer)
  formattedConfirm := confirmStyle.Render("Press Enter to confirm the commit message")

  if m.scope == "" {
    commitMessage := fmt.Sprintf("%s: %s %s\n\n%s\n\n%s\n\n%s",
      formattedCommitType, formattedGitmoji, formattedDescription, formattedBody, formattedFooter, formattedConfirm)
    return commitMessage
  }

  commitMessage := fmt.Sprintf("%s[%s]: %s %s\n\n%s\n\n%s\n\n%s",
    formattedCommitType, formattedScope, formattedGitmoji, formattedDescription, formattedBody, formattedFooter, formattedConfirm)

  return commitMessage
}

func (m model) CommitChanges() {
    // Create styles using lipgloss
    //headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Align(lipgloss.Left)
    successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Align(lipgloss.Left)
    //nextStepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Align(lipgloss.Left)

    // Create the git commit command
    cmd := exec.Command("git", "commit", "-m", m.PlainCommitMessage())
    cmd.Dir = "."

    // Run the command and capture the output
    _, err := cmd.Output()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Inform the user about the next steps with styles
    fmt.Println(successStyle.Render("Changes committed successfully. Now you can push the changes to the remote repository using: git push"))
    m.finished = true
}


func main() {
  cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")

  cmd.Dir = "."

  output, err := cmd.Output()

  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  result := strings.TrimSpace(string(output))

  if result != "true" {
    fmt.Println("Current directory is not a Git repository.")
    os.Exit(1)
  }

  p := tea.NewProgram(initialModel(), tea.WithAltScreen())

  if err := p.Start(); err != nil {
    fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
    os.Exit(1)
  }
}



