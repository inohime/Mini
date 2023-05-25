for i in {1..500}
    do curl -g "https://danbooru.donmai.us/tags.json?login={LOGIN}&api_key={API_KEY}&search[post_count.gte]=8000&search[order]=count&limit=5&page=$i" >> "v.json";
    if (( $i % 30 == 0 )); 
        then sleep 60; 
    fi
done

jq -s 'flatten' v.json > v_tmp.json && mv v_tmp.json v.json