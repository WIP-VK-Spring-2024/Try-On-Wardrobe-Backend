-- +migrate Up

-- +migrate StatementBegin
create function try_on_type(category text) returns text as $$
begin
    if category = 'Верх' then
        return 'upper_body';
    elsif category = 'Низ' then
        return 'lower_body';
    elsif category = 'Платья' then
        return 'dresses';
    else
        return '';
    end if;
end
$$ language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
drop function try_on_type;
