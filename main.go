package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Got unexpected error: %v", err)
	}
}

type Mode int

const (
	RssSelect Mode = iota
	ArticleSelect
	ArticleView
)

type model struct {
	mode Mode
	// article body in view mode
	// each elem is a line
	body []string

	articleStart int
	articleEnd   int

	rssItems []RSS

	rssCursor   int
	rssSelected int

	articleCursor   int
	articleSelected int

	err error
}

func initialModel() model {
	return model{
		articleStart: 0,
		// TODO: use some proper value like max height or smth
		articleEnd: 40,
	}
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type rssLoadSuccess struct{ rssItems []RSS }
type articleLoadSuccess struct{ body []string }

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		rssItems, err := fetchRssItems()

		if err != nil {
			return errMsg{err}
		}

		return rssLoadSuccess{rssItems}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		// is it usefull to store err in model?
		m.err = msg.err
		return m, tea.Quit
	case rssLoadSuccess:
		m.rssItems = msg.rssItems
		m.mode = RssSelect
		return m, nil
	case articleLoadSuccess:
		m.body = msg.body
		m.mode = ArticleView
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.mode == RssSelect && m.rssCursor > 0 {
				m.rssCursor--
			}
			if m.mode == ArticleSelect && m.articleCursor > 0 {
				m.articleCursor--
			}
			if m.mode == ArticleView && m.articleStart-10 >= 0 {
				m.articleStart -= 10
				m.articleEnd -= 10
			}
		case "down", "j":
			if m.mode == RssSelect && m.rssCursor < len(m.rssItems) {
				m.rssCursor++
			}
			if m.mode == ArticleSelect && m.articleCursor < len(m.rssItems[m.rssSelected].Channel.Items) {
				m.articleCursor++
			}
			if m.mode == ArticleView && m.articleEnd+10 <= len(m.body)-1 {
				m.articleStart += 10
				m.articleEnd += 10
			}
		case "enter":
			if m.mode == RssSelect {
				m.rssSelected = m.rssCursor
				m.mode = ArticleSelect
			} else if m.mode == ArticleSelect {
				m.articleSelected = m.articleCursor
				l := m.rssItems[m.rssSelected].Channel.Items[m.articleSelected].Link
				return m, func() tea.Msg {
					body, err := fetchArticle(l)

					if err != nil {
						return errMsg{err}
					}

					return articleLoadSuccess{body: strings.Split(body, "\n")}
				}
			}
		case "esc":
			if m.mode > 0 {
				m.mode--
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.mode == RssSelect {
		return renderRssSelect(m)
	}

	if m.mode == ArticleSelect {
		return renderArticleSelect(m)
	}

	if m.mode == ArticleView {
		return renderArticleView(m)
	}

	return "There will be list of articles here someday"
}

func renderRssSelect(m model) string {
	var sb strings.Builder

	for i, rss := range m.rssItems {

		cursor := " "
		if i == m.rssCursor {
			cursor = ">"
		}
		checked := " "
		if i == m.rssSelected {
			checked = "x"
		}
		s := fmt.Sprintf("%s [%s] %s\n", cursor, checked, rss.Channel.Title)
		sb.WriteString(s)
	}

	return sb.String()
}

func renderArticleSelect(m model) string {
	var sb strings.Builder

	for i, item := range m.rssItems[m.rssSelected].Channel.Items {
		cursor := " "
		if i == m.articleCursor {
			cursor = ">"
		}
		checked := " "
		if i == m.articleSelected {
			checked = "x"
		}
		s := fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.Title)
		sb.WriteString(s)
	}

	return sb.String()

}

func renderArticleView(m model) string {
	return strings.Join(m.body[m.articleStart:m.articleEnd], "\n")
}
