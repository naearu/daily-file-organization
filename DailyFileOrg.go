//go:generate rsrc -ico move_icon.ico -o rsrc.syso

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	fmt.Println("Daily File Organization Tool")

	// 현재 날짜 가져오기
	today := time.Now()
	todayStr := today.Format("20060102")

	// 현재 디렉토리에서 파일들의 날짜 추출하여 참조 날짜 구하기
	referenceDate := extractLatestDateFromFiles()
	if referenceDate == "" {
		// 파일에서 날짜를 찾을 수 없는 경우 기존 로직 사용
		if today.Weekday() == time.Monday {
			// 월요일이면 3일 전 날짜 (금요일)
			reference := today.AddDate(0, 0, -3)
			referenceDate = reference.Format("20060102")
		} else {
			// 월요일이 아니면 어제 날짜
			reference := today.AddDate(0, 0, -1)
			referenceDate = reference.Format("20060102")
		}
	}

	fmt.Printf("참조 날짜: %s\n", referenceDate)
	fmt.Printf("오늘 날짜: %s\n", todayStr)

	// 2025 폴더 생성
	if err := os.MkdirAll("2025", 0755); err != nil {
		log.Printf("폴더 생성 오류: %v", err)
	}

	// 현재 디렉토리에서 파일 찾기
	files, err := filepath.Glob("*")
	if err != nil {
		log.Fatalf("파일 검색 오류: %v", err)
	}

	// Excel 및 Word 파일 처리
	targetExtensions := []string{".xlsx", ".docx"}

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		isTargetFile := false

		for _, targetExt := range targetExtensions {
			if ext == targetExt {
				isTargetFile = true
				break
			}
		}

		if !isTargetFile {
			continue
		}

		// 파일을 2025 폴더로 복사
		if err := copyFile(file, filepath.Join("2025", file)); err != nil {
			log.Printf("파일 복사 오류 (%s): %v", file, err)
		} else {
			fmt.Printf("파일을 2025 폴더로 복사: %s\n", file)
		}

		// 파일명에 참조 날짜가 포함되어 있는지 확인하고 변경
		filename := strings.TrimSuffix(file, ext)
		if strings.Contains(filename, referenceDate) {
			newFilename := strings.ReplaceAll(filename, referenceDate, todayStr)
			newFile := newFilename + ext

			if err := os.Rename(file, newFile); err != nil {
				log.Printf("파일 이름 변경 오류 (%s): %v", file, err)
			} else {
				fmt.Printf("파일 이름 변경: %s → %s\n", file, newFile)
			}
		}
	}

	fmt.Println("파일 이름 변경 작업이 완료되었습니다.")
	fmt.Print("계속하려면 Enter 키를 누르세요...")
	fmt.Scanln()
}

// 파일 복사 함수
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 파일 권한 복사
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// 현재 디렉토리의 파일들에서 날짜를 추출하여 가장 최근 날짜를 반환
func extractLatestDateFromFiles() string {
	files, err := filepath.Glob("*")
	if err != nil {
		log.Printf("파일 검색 오류: %v", err)
		return ""
	}

	// 날짜 패턴 정규식: (YYYYMMDD) 형태
	datePattern := regexp.MustCompile(`\((\d{8})\)`)
	var latestDate string

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".xlsx" || ext == ".docx" {
			matches := datePattern.FindStringSubmatch(file)
			if len(matches) > 1 {
				dateStr := matches[1]
				// 날짜 형식이 유효한지 확인
				if _, err := time.Parse("20060102", dateStr); err == nil {
					if latestDate == "" || dateStr > latestDate {
						latestDate = dateStr
					}
				}
			}
		}
	}

	return latestDate
}
