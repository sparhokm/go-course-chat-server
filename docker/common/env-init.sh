#!/usr/bin/env bash

# Проверяем, что передано достаточно аргументов
if [ $# -lt 2 ]; then
    echo "Usage: $0 <output_file> <input_files...>"
    exit 1
fi

# Получаем имя файла, в который будем записывать результат
output_file=$1
shift

# Создаем пустой файл, чтобы не добавлять данные к существующему файлу
> "$output_file"

# Обходим все переданные файлы
for input_file in "$@"; do
    # Проверяем, что файл существует и доступен для чтения
    if [ ! -r "$input_file" ]; then
        continue
    fi

    # Читаем все строки из файла и добавляем их в выходной файл
    while read line; do
        # Ищем переменную в строке
        var_name=$(echo "$line" | cut -d= -f1)
        if grep -q "^$var_name=" "$output_file"; then
            # Если переменная уже есть в выходном файле, заменяем ее значение
            sed -i "s/^$var_name=.*/$line/" "$output_file"
        else
            # Иначе просто добавляем строку в выходной файл
            echo "$line" >> "$output_file"
        fi
    done < "$input_file"
done
