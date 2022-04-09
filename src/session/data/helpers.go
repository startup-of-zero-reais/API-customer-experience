package data

import "github.com/startup-of-zero-reais/API-customer-experience/src/session/domain"

type order string

const (
	ASC  = order("ASC")
	DESC = order("DESC")
)

func merge(firstPart, secondPart []domain.UserSession) []domain.UserSession {
	var n = make([]domain.UserSession, len(firstPart)+len(secondPart))

	var firstIndex = 0
	var secondIndex = 0

	var nIndex = 0

	for firstIndex < len(firstPart) && secondIndex < len(secondPart) {
		if firstPart[firstIndex].CreatedAt > secondPart[secondIndex].CreatedAt {
			n[nIndex] = firstPart[firstIndex]
			firstIndex++
		} else {
			n[nIndex] = secondPart[secondIndex]
			secondIndex++
		}

		nIndex++
	}

	for firstIndex < len(firstPart) {
		n[nIndex] = firstPart[firstIndex]
		firstIndex++
		nIndex++
	}

	for secondIndex < len(secondPart) {
		n[nIndex] = secondPart[secondIndex]
		secondIndex++
		nIndex++
	}

	return n
}

func sort(sessions []domain.UserSession) []domain.UserSession {
	if len(sessions) <= 1 {
		return sessions
	}

	var middle = len(sessions) / 2

	var left = sort(sessions[0:middle])
	var right = sort(sessions[middle:])

	return merge(left, right)
}
