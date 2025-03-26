package dna

// type MockMultipartFile struct {
// 	multipart.File
// 	mock.Mock
// 	data *bytes.Buffer
// }
//
// func (m *MockMultipartFile) Read(p []byte) (n int, err error) {
// 	n = copy(p, m.data.Next(len(p)))
// 	if m.data.Len() == 0 {
// 		return n, io.EOF
// 	}
// 	return n, nil
// }
//
// func (m *MockMultipartFile) Close() error {
// 	return nil
// }
//
// func TestAnalyzeGCContent(t *testing.T) {
// 	// Define test cases
// 	tests := []struct {
// 		name        string
// 		fileContent string
// 		expected    float64
// 	}{
// 		{
// 			name:        "Simple GC Content",
// 			fileContent: ">seq1\nGCGCGC\n>seq2\nCGCGCG",
// 			expected:    100,
// 		},
// 		{
// 			name:        "Mixed Content",
// 			fileContent: ">seq1\nGCGCAT\n>seq2\nATGCAT",
// 			expected:    50,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Mock the multipart file
// 			fileContent := tt.fileContent
// 			file := &MockMultipartFile{
// 				data: bytes.NewBufferString(fileContent),
// 			}
//
// 			service := NewDNAService()
// 			result, err := service.AnalyzeGCContent(file)
//
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }

	//
	// 	totalGCCount := totalCounts.G + totalCounts.C
	// 	totalSeqLength := totalGCCount + totalCounts.A + totalCounts.T

	// 	if totalSeqLength == 0 {
	// 		return 0, nil
	// 	}
	//
	// 	x := (float64(totalGCCount) / float64(totalSeqLength)) * 100
	// 	gcContent := utils.Round(x, gcPrecision)
