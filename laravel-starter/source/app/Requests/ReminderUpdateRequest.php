<?php

namespace App\Requests;

use App\Rules\ValidRRule;
use Illuminate\Foundation\Http\FormRequest;

class ReminderUpdateRequest extends FormRequest
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
            'rrule' => ['sometimes', 'string', new ValidRRule()],
            'description' => 'sometimes|string',
            'start_at' => 'sometimes|date',
        ];
    }
}
