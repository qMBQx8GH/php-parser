package printer_test

import (
	"testing"

	"github.com/Demooon86/php-parser/internal/tester"
)

func TestParseAndPrintPropertyHooksPHP8(t *testing.T) {
	tester.NewParserPrintTestSuite(t).UsePHP84().Run(`<?php
    class Example
    {
        public string $foo = 'default value' {
        get {
            return $this->foo;
        }

        set(string $value) {
            $this->foo = strtolower($value);
        }
    }
}
`)
}
