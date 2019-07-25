package k8s

import (
	"testing"
	"time"

	"gotest.tools/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInsertJobIntoSliceByCreationTime(t *testing.T) {
	a := Job{
		Name: "a",
		CreationTime: metav1.Time{
			time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	b := Job{
		Name: "b",
		CreationTime: metav1.Time{
			time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	c := Job{
		Name: "c",
		CreationTime: metav1.Time{
			time.Date(2008, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	type testCase struct {
		tName         string
		job           Job
		slice         []Job
		expectedSlice []Job
	}
	testCases := []testCase{
		testCase{
			tName:         "empty slice",
			job:           a,
			slice:         []Job{},
			expectedSlice: []Job{a},
		},
		testCase{
			tName:         "add to beginning",
			job:           a,
			slice:         []Job{b, c},
			expectedSlice: []Job{a, b, c},
		},
		testCase{
			tName:         "add to middle",
			job:           b,
			slice:         []Job{a, c},
			expectedSlice: []Job{a, b, c},
		},
		testCase{
			tName:         "add to end",
			job:           c,
			slice:         []Job{a, b},
			expectedSlice: []Job{a, b, c},
		},
		testCase{
			tName:         "equal to existing",
			job:           b,
			slice:         []Job{a, b, c},
			expectedSlice: []Job{a, b, b, c},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.tName, func(t *testing.T) {
			t.Parallel()
			result := insertJobIntoSliceByCreationTime(tc.slice, tc.job)
			assert.DeepEqual(t, tc.expectedSlice, result)
		})
	}
}

func TestInsertCronCronJobIntoSliceByCreationTime(t *testing.T) {
	a := CronJob{
		Name: "a",
		CreationTime: metav1.Time{
			time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	b := CronJob{
		Name: "b",
		CreationTime: metav1.Time{
			time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	c := CronJob{
		Name: "c",
		CreationTime: metav1.Time{
			time.Date(2008, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	type testCase struct {
		tName         string
		job           CronJob
		slice         []CronJob
		expectedSlice []CronJob
	}
	testCases := []testCase{
		testCase{
			tName:         "empty slice",
			job:           a,
			slice:         []CronJob{},
			expectedSlice: []CronJob{a},
		},
		testCase{
			tName:         "add to beginning",
			job:           a,
			slice:         []CronJob{b, c},
			expectedSlice: []CronJob{a, b, c},
		},
		testCase{
			tName:         "add to middle",
			job:           b,
			slice:         []CronJob{a, c},
			expectedSlice: []CronJob{a, b, c},
		},
		testCase{
			tName:         "add to end",
			job:           c,
			slice:         []CronJob{a, b},
			expectedSlice: []CronJob{a, b, c},
		},
		testCase{
			tName:         "equal to existing",
			job:           b,
			slice:         []CronJob{a, b, c},
			expectedSlice: []CronJob{a, b, b, c},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.tName, func(t *testing.T) {
			t.Parallel()
			result := insertCronJobIntoSliceByCreationTime(tc.slice, tc.job)
			assert.DeepEqual(t, tc.expectedSlice, result)
		})
	}
}
