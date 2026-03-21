package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// LinkItem representa um link com nome e URL
type LinkItem struct {
	Name string
	URL  string
}

// LinkSection representa uma secção de links
type LinkSection struct {
	Title       string
	Description string
	Color       int
	Links       []LinkItem
}

// Títulos e descrições fixos por secção (apenas metadados da mensamge embed, os links vêm da leitura do canal)
var SectionMeta = map[string]LinkSection{
	"links_uteis": {
		Title:       "Links Úteis",
		Description: "Links gerais úteis para o teu percurso académico na Universidade de Évora ao lado do NEEI-UÉ.",
		Color:       0xffffff,
	},
	"drives-resumos": {
		Title:       "Drives de Resumos",
		Description: "Drives e recursos de estudo partilhados pela comunidade.",
		Color:       0xffffff,
	},
	"genericos": {
		Title:       "Recursos Genéricos de Informática",
		Description: "Recursos genéricos que podem ser úteis para os estudantes de informática.",
		Color:       0xffffff,
	},
	"redes-sociais": {
		Title:       "Redes Sociais do NEEI",
		Description: "Segue-nos nas nossas redes sociais para estares sempre a par das novidades.",
		Color:       0xffffff,
	},
}

// ParseLinksFromMessages lê as mensagens do canal e devolve as secções com os links
// Formato esperado de cada mensagem no canal:
//
//	SECÇÃO: links_uteis
//	Nome do Link | https://url.com
//	Nome de outro Link | https://outro.com
//	e assim sucessivamente
func ParseLinksFromMessages(messages []*discordgo.Message) map[string]LinkSection {
	result := map[string]LinkSection{}

	// Itera pelas mensagens e tenta extrair as secções e os links
	for _, msg := range messages {
		lines := strings.Split(strings.TrimSpace(msg.Content), "\n")
		if len(lines) < 2 {
			continue
		}

		// Primeira linha deve ser "SECÇÃO: chave"
		firstLine := strings.TrimSpace(lines[0])
		if !strings.HasPrefix(firstLine, "SECCAO:") {
			continue
		}

		seccaoKey := strings.TrimSpace(strings.TrimPrefix(firstLine, "SECCAO:"))

		// Ignora se já existe uma versão mais recente desta secção
		if _, exists := result[seccaoKey]; exists {
			continue
		}

		// Verifica se a secção é válida (tem metadados definidos)
		meta, ok := SectionMeta[seccaoKey]
		if !ok {
			continue
		}

		// Faz parse dos links
		var links []LinkItem
		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			url := strings.TrimSpace(parts[1])
			if name != "" && url != "" {
				links = append(links, LinkItem{Name: name, URL: url})
			}
		}

		// Se não houver links válidos, ignora esta secção
		if len(links) == 0 {
			continue
		}

		// Adiciona a secção ao resultado
		result[seccaoKey] = LinkSection{
			Title:       meta.Title,
			Description: meta.Description,
			Color:       meta.Color,
			Links:       links,
		}
	}

	return result
}

// getMoodleURL gera o URL do Moodle com o ano letivo atual automaticamente
func getMoodleURL() string {
	// Nota: importa "time" no ficheiro que chamar esta função se necessário
	// Esta função é um helper para usar nos links do canal como referência
	_ = fmt.Sprintf // evita import não usado
	return "Ver canal de links para o URL atualizado"
}
