select inflation, income_growth, min(age) as age from yearly_stats
where can_cover = 1
  and inflation = 10
group by inflation, income_growth order by inflation, income_growth;


