<?php

namespace App\Requests;

use App\Rules\ValidRRule;
use Illuminate\Foundation\Http\FormRequest;

class ReminderCreateRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return auth()->check();
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, \Illuminate\Contracts\Validation\ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'rrule' => ['required', 'string', new ValidRRule()],
            'description' => 'required|string',
            'start_at' => 'required|date',
        ];
    }
}
