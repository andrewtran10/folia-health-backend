<?php

namespace App\Rules;

use RRule\RRule;
use DateTime;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class ValidRRule implements ValidationRule
{
    /**
     * Run the validation rule.
     *
     * @param  \Closure(string): \Illuminate\Translation\PotentiallyTranslatedString  $fail
     */
    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        if (!$this->isValidRRule($value)) {
            $fail("'$value' is not a valid RRULE string. Review RFC 5545 for valid formatting.");
        }
    }

    private function isValidRRule(string $rrule): bool
    {
        try {
            $rruleObject = new RRule($rrule);
        } catch (\InvalidArgumentException $e) {
            return false;
        }
        return true;
    }

}
