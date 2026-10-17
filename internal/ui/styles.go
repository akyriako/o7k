package ui

import (
	"github.com/akyriako/o7k/internal/config"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

//Header and AppView

var (
	resourceTagStyle          lipgloss.Style
	navigationTagStyle        lipgloss.Style
	commandStyle              lipgloss.Style
	commandGhostStyle         lipgloss.Style
	statusStyle               lipgloss.Style
	headerLabelStyle          lipgloss.Style
	headerValueStyle          lipgloss.Style
	headerCommandKeyStyle     lipgloss.Style
	headerCommandTextStyle    lipgloss.Style
	logoStyle                 lipgloss.Style
	loadingStyle              lipgloss.Style
	infoStyle                 lipgloss.Style
	autoRefreshStateOnStyle   lipgloss.Style
	autoRefreshStateOffStyle  lipgloss.Style
	modalButtonStyle          lipgloss.Style
	jsonKeyStyle              lipgloss.Style
	jsonStringStyle           lipgloss.Style
	jsonNumberStyle           lipgloss.Style
	jsonBoolStyle             lipgloss.Style
	jsonNullStyle             lipgloss.Style
	detailScrollbarTrackStyle lipgloss.Style
	detailScrollbarThumbStyle lipgloss.Style
	errorModalStyle           lipgloss.Style
	errorTitleStyle           lipgloss.Style
	tableTitleStyle           lipgloss.Style
	errorScrollbarTrackStyle  lipgloss.Style
	errorScrollbarThumbStyle  lipgloss.Style
)

func initStyles(skin config.Skin) {
	resourceTagStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ResourceTag.Foreground)).
		Background(lipgloss.Color(skin.ResourceTag.Background)).
		Bold(true)

	navigationTagStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.NavigationTag.Foreground)).
		Background(lipgloss.Color(skin.NavigationTag.Background))

	commandStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Command))

	commandGhostStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.CommandGhost)).
		Faint(true)

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Status.Foreground)).
		Background(lipgloss.Color(skin.Status.Background)).
		Align(lipgloss.Center)

	headerLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderLabel)).
		Bold(true)

	headerValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderValue))

	headerCommandKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderCommandKey)).
		Bold(true)

	headerCommandTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderCommandText))

	logoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Logo)).
		Bold(true)

	loadingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Loading.Foreground)).
		Background(lipgloss.Color(skin.Loading.Background)).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Info.Foreground)).
		Background(lipgloss.Color(skin.Info.Background)).
		Bold(true)

	autoRefreshStateOnStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.StateOn.Foreground)).
		Background(lipgloss.Color(skin.StateOn.Background)).
		Bold(true)

	autoRefreshStateOffStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.StateOff.Foreground)).
		Background(lipgloss.Color(skin.StateOff.Background)).
		Bold(true)

	modalButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ModalButton.Foreground)).
		Background(lipgloss.Color(skin.ModalButton.Background)).
		Bold(true)

	//DetailsView

	jsonKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONKey))

	jsonStringStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONString))

	jsonNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONNumber))

	jsonBoolStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONBool))

	jsonNullStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONNull))

	detailScrollbarTrackStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.DetailBorder))

	detailScrollbarThumbStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.DetailBorder)).
		Bold(true)

	//ErrorView

	errorModalStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(skin.ErrorBorder)).
		Foreground(lipgloss.Color(skin.ErrorModal.Foreground)).
		Padding(1, 1).
		Width(60)

	errorTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ErrorTitle)).
		Bold(true)

	//TableView

	tableTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.TableTitle)).Bold(true)

	errorScrollbarTrackStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ErrorBorder))

	errorScrollbarThumbStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ErrorBorder))

	if skin.ErrorModal.Background != "" {
		errorModalStyle = errorModalStyle.Background(lipgloss.Color(skin.ErrorModal.Background))
	}
}

func tableStyles(skin config.Skin) table.Styles {
	styles := table.DefaultStyles()

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color(skin.TableSelected.Foreground)).
		Background(lipgloss.Color(skin.TableSelected.Background)).
		Bold(true)

	return styles
}
