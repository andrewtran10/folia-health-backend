<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Reminder extends Model
{
    use HasFactory, HasUuids;

    protected $fillable = ['user_id', 'rrule', 'description', 'start_at'];

    protected function casts(): array
    {
        return [
            'start_at' => 'datetime',
        ];
    }

    // Define reminder -> user relationship
    public function user()
    {
        return $this->belongsTo(User::class);
    }
}
