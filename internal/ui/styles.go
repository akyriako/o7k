package ui

import (
	"github.com/akyriako/o7k/internal/config"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

//Header and AppView

var (
	mainStyle                 lipgloss.Style
	tableBodyStyle            lipgloss.Style
	detailBodyStyle           lipgloss.Style
	resourceTagStyle          lipgloss.Style
	navigationTagStyle        lipgloss.Style
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
	mainStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Main.Foreground)).
		Background(lipgloss.Color(skin.Main.Background))

	tableBodyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(skin.Table.Foreground)).Background(lipgloss.Color(skin.Table.Background))
	detailBodyStyle = mainStyle

	resourceTagStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ResourceTag.Foreground)).
		Background(lipgloss.Color(skin.ResourceTag.Background)).
		Bold(true)

	navigationTagStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.NavigationTag.Foreground)).
		Background(lipgloss.Color(skin.NavigationTag.Background))

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Status.Foreground)).
		Background(lipgloss.Color(skin.Status.Background)).
		Align(lipgloss.Center)

	headerLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderLabel)).
		Background(lipgloss.Color(skin.Context.Background)).
		Bold(true)

	headerValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderValue)).
		Background(lipgloss.Color(skin.Context.Background))

	headerCommandKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderCommandKey)).
		Background(lipgloss.Color(skin.Main.Background)).
		Bold(true)

	headerCommandTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.HeaderCommandText)).
		Background(lipgloss.Color(skin.Main.Background))

	logoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Logo)).
		Background(lipgloss.Color(skin.Main.Background)).
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
		Foreground(lipgloss.Color(skin.JSONKey)).
		Background(lipgloss.Color(skin.Main.Background))

	jsonStringStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONString)).
		Background(lipgloss.Color(skin.Main.Background))

	jsonNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONNumber)).
		Background(lipgloss.Color(skin.Main.Background))

	jsonBoolStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONBool)).
		Background(lipgloss.Color(skin.Main.Background))

	jsonNullStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.JSONNull)).
		Background(lipgloss.Color(skin.Main.Background))

	detailScrollbarTrackStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Main.Foreground)).
		Background(lipgloss.Color(skin.Main.Background))

	detailScrollbarThumbStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Main.Foreground)).
		Background(lipgloss.Color(skin.Main.Background)).
		Bold(true)

	//ErrorView

	errorModalStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(skin.Main.Foreground)).
		BorderBackground(lipgloss.Color(skin.Main.Background)).
		Foreground(lipgloss.Color(skin.Main.Foreground)).
		Background(lipgloss.Color(skin.Main.Background)).
		Padding(1, 1).
		Width(60)

	errorTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.ErrorTitle)).
		Bold(true)

	errorScrollbarTrackStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Main.Foreground))

	errorScrollbarThumbStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.Main.Foreground))

	//TableView

	tableTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(skin.TableTitle)).Bold(true)
}

func tableStyles(skin config.Skin) table.Styles {
	styles := table.DefaultStyles()

	styles.Cell = styles.Cell.
		Foreground(lipgloss.Color(skin.Table.Foreground)).
		Background(lipgloss.Color(skin.Table.Background))

	styles.Header = styles.Header.
		Foreground(lipgloss.Color(skin.TableHeader)).
		Background(lipgloss.Color(skin.Table.Background))

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color(skin.TableSelected.Foreground)).
		Background(lipgloss.Color(skin.TableSelected.Background)).
		Bold(true)

	return styles
}
